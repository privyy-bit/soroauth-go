package readmesnippets

import (
	"context"

	"github.com/stellar/go-stellar-sdk/keypair"
	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/soroauth/soroauth-go"
)

// SignBatchEntries backs the "Sign a batch of authorization entries" recipe.
// It signs multiple entries with the same signer and expiration, returning all
// signed entries or none on error (all-or-nothing, per AuthorizeAll).
// See: docs/batch.go#AuthorizeAll
func SignBatchEntries(
	ctx context.Context,
	entries []xdr.SorobanAuthorizationEntry,
	signer *keypair.Full,
	validUntil uint32,
	passphrase string,
) ([]xdr.SorobanAuthorizationEntry, error) {
	// snippet:start sign_batch
	signed, err := soroauth.AuthorizeAll(ctx, entries,
		[]soroauth.Signer{soroauth.NewEd25519Signer(signer)},
		validUntil, passphrase)
	// snippet:end sign_batch
	if err != nil {
		return nil, err // nothing partial is ever returned
	}
	return signed, nil
}
