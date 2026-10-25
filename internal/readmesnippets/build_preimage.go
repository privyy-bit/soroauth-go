package readmesnippets

import (
	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/soroauth/soroauth-go"
)

// BuildPreimageAndPayload backs the "Build a preimage and payload hash"
// recipe. It computes the HashIdPreimage from an entry, validUntilLedger, and
// network passphrase, then returns the preimage XDR and its SHA-256 payload.
// The payload is the 32-byte hash used as the basis for signatures.
// See: docs/preimage.go#Preimage
func BuildPreimageAndPayload(
	entry xdr.SorobanAuthorizationEntry,
	validUntilLedger uint32,
	networkPassphrase string,
) (xdr.HashIdPreimage, [32]byte, error) {
	// snippet:start build_preimage
	preimage, err := soroauth.Preimage(entry, validUntilLedger, networkPassphrase)
	if err != nil {
		return xdr.HashIdPreimage{}, [32]byte{}, err
	}
	payload, err := soroauth.Payload(preimage)
	if err != nil {
		return xdr.HashIdPreimage{}, [32]byte{}, err
	}
	// snippet:end build_preimage
	return preimage, payload, nil
}
