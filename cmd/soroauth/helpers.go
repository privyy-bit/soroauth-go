package main

import (
	"io"
	"os"
)

func readAssertionInput(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func resolveEntryArg(val string, r io.Reader) (string, error) {
	if val == "-" {
		if r == nil {
			r = os.Stdin
		}
		data, err := io.ReadAll(r)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return val, nil
}
