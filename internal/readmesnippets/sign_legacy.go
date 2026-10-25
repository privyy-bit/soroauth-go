package readmesnippets

import (
	"context"

	"github.com/stellar/go-stellar-sdk/keypair"
	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/soroauth/soroauth-go"
)

// SignLegacyEntry backs the "Sign a single legacy address entry" recipe.
// It signs a legacy SOROBAN_CREDENTIALS_ADDRESS (arm 1) entry with an ed25519
// signer. The entry is assumed to be unsigned and from a simulation.
// See: docs/credential-types.md#SOROBAN_CREDENTIALS_ADDRESS
func SignLegacyEntry(
	ctx context.Context,
	entry xdr.SorobanAuthorizationEntry,
	signer *keypair.Full,
	validUntil uint32,
	passphrase string,
) (xdr.SorobanAuthorizationEntry, error) {
	// snippet:start sign_legacy
	signed, err := soroauth.AuthorizeEntry(ctx, entry, soroauth.NewEd25519Signer(signer), validUntil, passphrase)
	// snippet:end sign_legacy
	if err != nil {
		return xdr.SorobanAuthorizationEntry{}, err
	}
	return signed, nil
}
