package readmesnippets

import (
	"context"

	"github.com/stellar/go-stellar-sdk/keypair"
	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/soroauth/soroauth-go"
)

// SignDelegatesEntry backs the "Sign a delegates tree entry" recipe.
// It wraps an entry with delegates and signs each node (account + every
// delegate at every depth) for a V2 arm. The top-level signature may be
// ScvVoid if only delegates authenticate, per CAP-71-01.
// See: docs/delegates.md
func SignDelegatesEntry(
	ctx context.Context,
	entry xdr.SorobanAuthorizationEntry,
	validUntil uint32,
	passphrase string,
	d1, d2, d3 string,
	k1, k2, k3 *keypair.Full,
) (xdr.SorobanAuthorizationEntry, error) {
	// snippet:start sign_delegates
	wrapped, err := soroauth.WithDelegates(entry, validUntil,
		[]soroauth.Delegate{
			{Address: d1},
			{Address: d2, Nested: []soroauth.Delegate{{Address: d3}}},
		}, nil) // nil top-level signature → ScvVoid, which CAP-71-01 permits
	if err != nil {
		return xdr.SorobanAuthorizationEntry{}, err
	}
	for _, kp := range []*keypair.Full{k1, k2, k3} {
		wrapped, err = soroauth.AuthorizeEntry(ctx, wrapped,
			soroauth.NewEd25519Signer(kp), validUntil, passphrase,
			soroauth.ForAddress(kp.Address()))
		if err != nil {
			return xdr.SorobanAuthorizationEntry{}, err
		}
	}
	// snippet:end sign_delegates
	return wrapped, nil
}
