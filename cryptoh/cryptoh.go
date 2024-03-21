// Package cryptoh provides helper functions for the crypto package.
package cryptoh

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// Sha256FileHash computes the sha256 of the input file.
func Sha256FileHash(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	out := h.Sum(nil)
	if len(out) != 32 {
		return nil, fmt.Errorf("hash is incorrect length")
	}

	return out, nil
}
