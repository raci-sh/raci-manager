// Package canonical produces the canonical JSON form used throughout RACI.sh:
// RFC 8785 (JSON Canonicalization Scheme). Every digest in the log and in the
// .raci file is computed over this form (SPEC.md section 4.1).
package canonical

import (
	"encoding/json"
	"fmt"

	"github.com/gowebpki/jcs"
)

// Marshal encodes v as JSON and canonicalizes it per RFC 8785.
func Marshal(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("canonical: marshal: %w", err)
	}
	return Transform(raw)
}

// Transform canonicalizes an existing JSON document per RFC 8785.
func Transform(raw []byte) ([]byte, error) {
	out, err := jcs.Transform(raw)
	if err != nil {
		return nil, fmt.Errorf("canonical: transform: %w", err)
	}
	return out, nil
}
