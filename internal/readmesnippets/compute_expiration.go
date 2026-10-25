package readmesnippets

import (
	"fmt"

	"github.com/soroauth/soroauth-go"
)

// ComputeExpiration backs the "Compute expiration ledger" recipe.
// It returns the absolute ledger at which the signature expires: latestLedger +
// ledgers. Returns ErrInvalidExpiration if ledgers is 0 or the sum overflows.
// See: docs/expiration.go#ExpirationAfter
func ComputeExpiration(latestLedger, ledgers uint32) (uint32, error) {
	// snippet:start compute_expiration
	result, err := soroauth.ExpirationAfter(latestLedger, ledgers)
	// snippet:end compute_expiration
	if err != nil {
		return 0, err
	}
	fmt.Printf("Expiration ledger: %d\n", result)
	return result, nil
}
