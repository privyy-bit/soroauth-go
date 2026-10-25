package readmesnippets

import (
	"fmt"

	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/soroauth/soroauth-go"
)

// WrapLegacyWithDelegates backs the "Wrap a legacy entry with delegates" recipe.
// It converts a legacy arm entry into the delegates arm using WithDelegates,
// setting the top-level signature to ScvVoid (CAP-71-01 permits this when every
// delegate authenticates).
// See: docs/delegates.go#WithDelegates
func WrapLegacyWithDelegates(
	entry xdr.SorobanAuthorizationEntry,
	validUntil uint32,
	delegates []soroauth.Delegate,
) (xdr.SorobanAuthorizationEntry, error) {
	// snippet:start wrap_legacy
	result, err := soroauth.WithDelegates(entry, validUntil, delegates, nil)
	// snippet:end wrap_legacy
	if err != nil {
		return xdr.SorobanAuthorizationEntry{}, err
	}
	fmt.Printf("Wrapped entry with %d delegates\n", len(delegates))
	return result, nil
}
