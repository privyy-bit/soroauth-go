package soroauth

import (
	"bytes"
	"fmt"

	"github.com/stellar/go-stellar-sdk/xdr"
)

// Equal returns true if two SorobanAuthorizationEntry values are XDR byte-
// identical. The comparison is a flat binary comparison of the encoded XDR,
// so two entries that differ only in delegate order, nonce value, or expiration
// ledger are considered unequal even though they may authorize the same
// operation. This is intentionally a byte-equality check, not a semantic
// equivalence check.
//
// An encoding failure (e.g. an entry that cannot be marshalled) returns false.
// It is never wrapped in an error; callers should not treat a false result as
// a refusal or fault.
func Equal(a, b xdr.SorobanAuthorizationEntry) bool {
	ba, err := a.MarshalBinary()
	if err != nil {
		return false
	}
	bb, err := b.MarshalBinary()
	if err != nil {
		return false
	}
	return bytesEqual(ba, bb)
}

// Clone returns a deep copy of entry that shares no memory with the input.
// The copy is produced by round-tripping through XDR encoding, the same method
// used internally by internal/xdrcopy. The returned value is safe to mutate
// without affecting the original.
// is decoded internally via SafeUnmarshalBase64 for deep-copy purposes, not from
// untrusted caller input.
//
//go:decode-allow: the base64 blob is an XDR-encoded SorobanAuthorizationEntry that
func Clone(entry xdr.SorobanAuthorizationEntry) (xdr.SorobanAuthorizationEntry, error) {
	// Marshal the entry to XDR bytes.
	b, err := entry.MarshalBinary()
	if err != nil {
		return xdr.SorobanAuthorizationEntry{}, fmt.Errorf(
			"soroauth: clone: marshal: %w", err)
	}

	// Unmarshal back into a new entry. This allocates fresh memory for every
	// pointer and slice, yielding a deep copy that shares no storage with the
	// original.
	var copied xdr.SorobanAuthorizationEntry
	if err := xdr.SafeUnmarshalBase64(string(b), &copied); err != nil {
		return xdr.SorobanAuthorizationEntry{}, fmt.Errorf(
			"soroauth: clone: unmarshal: %w", err)
	}

	return copied, nil
}

// bytesEqual returns true if two byte slices are equal.
func bytesEqual(a, b []byte) bool {
	return len(a) == len(b) && bytes.Equal(a, b)
}
