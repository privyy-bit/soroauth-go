package readmesnippets

import (
	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)

// UpgradeLegacyToV2 backs the "Upgrade a legacy unsigned entry to V2" recipe.
// It converts an unsigned legacy arm entry to the V2 arm. If the entry already
// carries a non-Void signature, it returns ErrAlreadySigned because the old
// signature would be invalid under the new payload.
// See: docs/upgrade.go#UpgradeToV2
func UpgradeLegacyToV2(entry xdr.SorobanAuthorizationEntry) (xdr.SorobanAuthorizationEntry, error) {
	// snippet:start upgrade_to_v2
	result, err := soroauth.UpgradeToV2(entry)
	// snippet:end upgrade_to_v2
	if err != nil {
		return xdr.SorobanAuthorizationEntry{}, err
	}
	return result, nil
}
