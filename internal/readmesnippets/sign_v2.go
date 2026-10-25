package readmesnippets

import (
	"context"

	"github.com/stellar/go-stellar-sdk/keypair"
	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/soroauth/soroauth-go"
)

// SignV2Entry backs the "Sign a single V2 address-bound entry" recipe.
// It signs a V2 SOROBAN_CREDENTIALS_ADDRESS_V2 (arm 2) entry with an ed25519
// signer. The address is bound into the signed payload.
// See: docs/credential-types.md#SOROBAN_CREDENTIALS_ADDRESS_V2
func SignV2Entry(
	ctx context.Context,
	entry xdr.SorobanAuthorizationEntry,
	signer *keypair.Full,
	validUntil uint32,
	passphrase string,
) (xdr.SorobanAuthorizationEntry, error) {
	// snippet:start sign_v2
	signed, err := soroauth.AuthorizeEntry(ctx, entry, soroauth.NewEd25519Signer(signer), validUntil, passphrase)
	// snippet:end sign_v2
	if err != nil {
		return xdr.SorobanAuthorizationEntry{}, err
	}
	return signed, nil
}
