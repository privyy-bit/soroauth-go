package readmesnippets

import (
	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)

// InspectEntry backs the "Inspect an authorization entry" recipe.
// It reports the entry's credential arm, address, nonce, expiration, root
// invocation, delegate tree, and which nodes are signed.
// See: docs/inspect.go#Inspect
func InspectEntry(entry xdr.SorobanAuthorizationEntry) (soroauth.EntryInfo, error) {
	// snippet:start inspect_entry
	info, err := soroauth.Inspect(entry)
	// snippet:end inspect_entry
	if err != nil {
		return soroauth.EntryInfo{}, err
	}
	return info, nil
}
