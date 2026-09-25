package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/soroauth/soroauth-go"
)

const delegatesUsage = `soroauth delegates — wrap an entry in a delegated-signer credential.

usage:
  soroauth delegates --entry <base64> --valid-until <ledger> \
                     --delegate <address> [--delegate <address> ...] \
                     [--nested PARENT=CHILD ...] [--json]

Converts an ADDRESS or ADDRESS_V2 entry into ADDRESS_WITH_DELEGATES (CAP-71-01),
with the delegates sorted into the order the protocol requires. Pass --delegate
once per address; the order they are given in does not matter. Nested delegates
are specified with --nested PARENT=CHILD, where PARENT is an existing top-level
or nested delegate address.

The delegate signatures are left as placeholders. Fill each one afterwards with:

  soroauth sign --entry <wrapped> --for <delegate address> ...

Nested delegates are supported via --nested PARENT=CHILD flags, matching the
library's delegate tree model (soroauth.Delegate.Nested).

Note that wrapping a legacy ADDRESS entry makes its payload address-bound, so
any signature already on the entry would stop verifying; such an entry is
rejected rather than silently rewrapped.

Prints the wrapped entry as base64. With --json, prints a JSON object with field
"wrapped_entry". On error, prints a JSON object with field "error" to stdout and
exits non-zero.
`

type delegatesOutput struct {
	WrappedEntry string `json:"wrapped_entry,omitempty"`
	Error        string `json:"error,omitempty"`
}

// addressList collects a flag that may be repeated.
type addressList []string

// nestedList collects a --nested PARENT=CHILD flag that may be repeated.
type nestedList []string

func (n *nestedList) String() string { return fmt.Sprint(*n) }

func (n *nestedList) Set(value string) error {
	if value == "" {
		return fmt.Errorf("--nested needs PARENT=CHILD")
	}
	*n = append(*n, value)
	return nil
}

func (a *addressList) String() string { return fmt.Sprint(*a) }

func (a *addressList) Set(value string) error {
	if value == "" {
		return fmt.Errorf("--delegate needs an address")
	}
	*a = append(*a, value)
	return nil
}

func runDelegates(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("delegates", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		fmt.Fprint(stderr, delegatesUsage)
		fmt.Fprintln(stderr, "\nflags:")
		flags.PrintDefaults()
	}

	entryFlag := flags.String("entry", "", "the authorization entry, as base64 XDR")
	validUntil := flags.Uint("valid-until", 0, "the last ledger at which the signatures are valid")
	var delegates addressList
	flags.Var(&delegates, "delegate", "a delegate address; repeat for several")
	var nestedFlags nestedList
	flags.Var(&nestedFlags, "nested", "a nested delegate as PARENT=CHILD; repeat for several")
	jsonFlag := flags.Bool("json", false, "output as JSON")

	if err := flags.Parse(args); err != nil {
		return newErrorf(ExitUsageError, "%w", err)
	}

	entry, err := decodeEntry(*entryFlag)
	if err != nil {
		return writeJSONError(stdout, *jsonFlag, err)
	}
	if *validUntil == 0 {
		return writeJSONError(stdout, *jsonFlag, newErrorf(ExitUsageError, "--valid-until is required and must be greater than zero"))
	}
	if len(delegates) == 0 {
		return writeJSONError(stdout, *jsonFlag, newErrorf(ExitUsageError, "at least one --delegate is required"))
	}

	delegateMap := make(map[string]*soroauth.Delegate)
	var topLevel []soroauth.Delegate

	for _, address := range delegates {
		d := &soroauth.Delegate{Address: address}
		delegateMap[address] = d
		topLevel = append(topLevel, *d)
	}

	// If there are top-level delegates, make sure every delegate in delegateMap points to its actual entry in the slice if needed,
	// or build recursively.
	var findAndAddNested func(nodes []soroauth.Delegate, parent, child string) bool
	findAndAddNested = func(nodes []soroauth.Delegate, parent, child string) bool {
		for i := range nodes {
			if nodes[i].Address == parent {
				cd := soroauth.Delegate{Address: child}
				nodes[i].Nested = append(nodes[i].Nested, cd)
				delegateMap[child] = &nodes[i].Nested[len(nodes[i].Nested)-1]
				return true
			}
			if findAndAddNested(nodes[i].Nested, parent, child) {
				return true
			}
		}
		return false
	}

	for _, pair := range nestedFlags {
		var parent, child string
		var parsedOK bool
		for i := 0; i < len(pair); i++ {
			if pair[i] == '=' {
				parent = pair[:i]
				child = pair[i+1:]
				parsedOK = true
				break
			}
		}
		if !parsedOK || parent == "" || child == "" {
			return writeJSONError(stdout, *jsonFlag, newErrorf(ExitUsageError, "malformed --nested argument %q, expected PARENT=CHILD", pair))
		}
		// Ensure parent exists
		if _, ok := delegateMap[parent]; !ok {
			return writeJSONError(stdout, *jsonFlag, newErrorf(ExitUsageError, "parent delegate %q in --nested not found among delegates", parent))
		}
		// Add child
		added := false
		for i := range topLevel {
			if topLevel[i].Address == parent {
				cd := soroauth.Delegate{Address: child}
				topLevel[i].Nested = append(topLevel[i].Nested, cd)
				delegateMap[child] = &topLevel[i].Nested[len(topLevel[i].Nested)-1]
				added = true
				break
			}
			if findAndAddNested(topLevel[i].Nested, parent, child) {
				added = true
				break
			}
		}
		if !added {
			return writeJSONError(stdout, *jsonFlag, newErrorf(ExitUsageError, "parent delegate %q in --nested not found", parent))
		}
	}

	wrapped, err := soroauth.WithDelegates(entry, uint32(*validUntil), topLevel, nil)
	if err != nil {
		// Classify the error for exit code
		var exitCode int
		if errors.Is(err, soroauth.ErrUnsupportedCredentials) ||
			errors.Is(err, soroauth.ErrAlreadySigned) ||
			errors.Is(err, soroauth.ErrDuplicateDelegate) {
			exitCode = ExitSigningRefusal
		} else if errors.Is(err, soroauth.ErrInvalidExpiration) {
			exitCode = ExitVerificationFailed
		} else {
			exitCode = ExitGeneralError
		}
		return writeJSONError(stdout, *jsonFlag, newErrorf(exitCode, "%w", err))
	}

	encoded, err := encodeEntry(wrapped)
	if err != nil {
		return writeJSONError(stdout, *jsonFlag, newErrorf(ExitGeneralError, "%w", err))
	}

	if *jsonFlag {
		out := delegatesOutput{WrappedEntry: encoded}
		enc := json.NewEncoder(stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(out)
	}

	fmt.Fprintln(stdout, encoded)
	return nil
}
