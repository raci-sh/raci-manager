// Package digest implements the digest format of SPEC.md section 4.1: the
// string "blake3:" followed by 64 lowercase hexadecimal characters.
package digest

import (
	"encoding/hex"
	"fmt"
	"strings"

	"lukechampine.com/blake3"
)

// Prefix starts every digest.
const Prefix = "blake3:"

// Digest is a blake3 digest in its "blake3:<hex>" string form.
type Digest string

// Zero is the all-zero digest: the prev_hash of a realm's first event and the
// head_hash of an empty log.
var Zero = Digest(Prefix + strings.Repeat("0", 64))

// Sum returns the blake3-256 digest of b.
func Sum(b []byte) Digest {
	sum := blake3.Sum256(b)
	return Digest(Prefix + hex.EncodeToString(sum[:]))
}

// Parse validates s and returns it as a Digest.
func Parse(s string) (Digest, error) {
	hexPart, ok := strings.CutPrefix(s, Prefix)
	if !ok {
		return "", fmt.Errorf("digest %q: missing %q prefix", s, Prefix)
	}
	if len(hexPart) != 64 {
		return "", fmt.Errorf("digest %q: want 64 hex characters, got %d", s, len(hexPart))
	}
	for _, c := range hexPart {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return "", fmt.Errorf("digest %q: not lowercase hexadecimal", s)
		}
	}
	return Digest(s), nil
}

// Valid reports whether d is a well-formed digest.
func (d Digest) Valid() bool {
	_, err := Parse(string(d))
	return err == nil
}
