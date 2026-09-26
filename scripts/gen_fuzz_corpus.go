// Command gen_fuzz_corpus builds fuzz seed corpora from the committed golden vectors.
//
// Usage:
//   go run ./scripts/gen_fuzz_corpus.go
package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type vectorData struct {
	Name             string `json:"name"`
	PreWrapEntryXDR  string `json:"pre_wrap_entry_xdr"`
	UnsignedEntryXDR string `json:"unsigned_entry_xdr"`
	SignedEntryXDR   string `json:"signed_entry_xdr"`
	PreimageXDR      string `json:"preimage_xdr"`
	PayloadHex       string `json:"payload_hex"`
}

func main() {
	vectorsDir := filepath.Join("testdata", "vectors")
	entries, err := os.ReadDir(vectorsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading vectors dir: %v\n", err)
		os.Exit(1)
	}

	var vectors []vectorData
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		path := filepath.Join(vectorsDir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading %s: %v\n", path, err)
			os.Exit(1)
		}
		var v vectorData
		if err := json.Unmarshal(raw, &v); err != nil {
			fmt.Fprintf(os.Stderr, "decoding %s: %v\n", path, err)
			os.Exit(1)
		}
		vectors = append(vectors, v)
	}

	// Define fuzz target corpus directories (under testdata/fuzz/)
	// Adjust or add corpus directories corresponding to fuzz targets in the repo.
	targets := []string{
		"FuzzDecodeEntry",
		"FuzzAuthorizeEntry",
		"FuzzVerifyEntry",
	}

	for _, target := range targets {
		corpusDir := filepath.Join("testdata", "fuzz", target)
		if err := os.MkdirAll(corpusDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "creating corpus dir %s: %v\n", corpusDir, err)
			os.Exit(1)
		}

		for _, v := range vectors {
			var payload []byte
			if v.UnsignedEntryXDR != "" {
				if b, err := base64.StdEncoding.DecodeString(v.UnsignedEntryXDR); err == nil {
					payload = append(payload, b...)
				}
			}
			if v.PreWrapEntryXDR != "" {
				if b, err := base64.StdEncoding.DecodeString(v.PreWrapEntryXDR); err == nil {
					payload = append(payload, b...)
				}
			}
			if v.SignedEntryXDR != "" {
				if b, err := base64.StdEncoding.DecodeString(v.SignedEntryXDR); err == nil {
					payload = append(payload, b...)
				}
			}
			if len(payload) == 0 && v.PayloadHex != "" {
				if b, err := hex.DecodeString(v.PayloadHex); err == nil {
					payload = b
				}
			}
			if len(payload) == 0 {
				continue
			}

			hash := sha256.Sum256(payload)
			fileName := fmt.Sprintf("%s_%x", v.Name, hash[:8])
			outPath := filepath.Join(corpusDir, fileName)
			if err := os.WriteFile(outPath, payload, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "writing corpus seed %s: %v\n", outPath, err)
				os.Exit(1)
			}
		}
	}
	fmt.Println("Fuzz seed corpus successfully generated from golden vectors.")
}
