// Package main implements a corpus generation utility that extracts XDR bytes
// from the committed golden vectors and populates fuzz seed directories.
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type vectorEntry struct {
	Name             string `json:"name"`
	PreWrapEntryXDR  string `json:"pre_wrap_entry_xdr"`
	UnsignedEntryXDR string `json:"unsigned_entry_xdr"`
	PreimageXDR      string `json:"preimage_xdr"`
	SignedEntryXDR   string `json:"signed_entry_xdr"`
}

func main() {
	vectorsDir := filepath.Join("testdata", "vectors")
	corpusDir := filepath.Join("testdata", "corpus")

	entries, err := os.ReadDir(vectorsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gen_corpus: reading vectors dir: %v\n", err)
		os.Exit(1)
	}

	targets := []string{"FuzzAuthorizeEntry", "FuzzPreimage", "FuzzPayload", "FuzzInspect"}
	for _, target := range targets {
		if err := os.MkdirAll(filepath.Join(corpusDir, target), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "gen_corpus: creating corpus dir for %s: %v\n", target, err)
			os.Exit(1)
		}
	}

	vectorCount := 0
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		path := filepath.Join(vectorsDir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gen_corpus: reading %s: %v\n", path, err)
			os.Exit(1)
		}

		var v vectorEntry
		if err := json.Unmarshal(raw, &v); err != nil {
			fmt.Fprintf(os.Stderr, "gen_corpus: decoding %s: %v\n", path, err)
			os.Exit(1)
		}

		addSeed := func(target, xdrStr string, suffix string) {
			if xdrStr == "" {
				return
			}
			bs, err := base64.StdEncoding.DecodeString(xdrStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "gen_corpus: base64 decode error in %s (%s): %v\n", v.Name, suffix, err)
				return
			}
			outPath := filepath.Join(corpusDir, target, fmt.Sprintf("%s_%s.bin", v.Name, suffix))
			if err := os.WriteFile(outPath, bs, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "gen_corpus: writing seed %s: %v\n", outPath, err)
				os.Exit(1)
			}
		}

		addSeed("FuzzAuthorizeEntry", v.UnsignedEntryXDR, "unsigned")
		addSeed("FuzzAuthorizeEntry", v.PreWrapEntryXDR, "prewrap")
		addSeed("FuzzPreimage", v.PreimageXDR, "preimage")
		addSeed("FuzzPayload", v.PreimageXDR, "preimage")
		addSeed("FuzzInspect", v.SignedEntryXDR, "signed")
		addSeed(
			"FuzzInspect", v.UnsignedEntryXDR, "unsigned")

		vectorCount++
	}

	fmt.Printf("Successfully generated fuzz corpus from %d golden vectors.\n", vectorCount)
}
