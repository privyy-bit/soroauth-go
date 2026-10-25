package readmesnippets

import (
	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)

// ValidateDelegateOrdering backs the "Validate delegate ordering" recipe.
// It checks an existing delegates-arm entry for correct sorting (ascending by
// XDR bytes of address) and for duplicate addresses within each level.
// See: docs/delegates.go#ValidateDelegateOrder
func ValidateDelegateOrdering(entry xdr.SorobanAuthorizationEntry) error {
	// snippet:start validate_delegates
	err := soroauth.ValidateDelegateOrder(entry)
	// snippet:end validate_delegates
	return err
}
