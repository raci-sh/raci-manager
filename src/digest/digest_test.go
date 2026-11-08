package digest

import "testing"

func TestZero(t *testing.T) {
	if len(Zero) != len(Prefix)+64 || !Zero.Valid() {
		t.Errorf("Zero = %q is not a valid digest", Zero)
	}
}

func TestSumFormat(t *testing.T) {
	d := Sum([]byte("raci"))
	if !d.Valid() || d == Zero {
		t.Errorf("Sum = %q", d)
	}
	if Sum([]byte("raci")) != d {
		t.Error("Sum is not deterministic")
	}
}

func TestParseRejects(t *testing.T) {
	for _, s := range []string{
		"sha256:" + string(Zero[len(Prefix):]),
		Prefix + "00",
		Prefix + "ABCDEF0000000000000000000000000000000000000000000000000000000000",
	} {
		if _, err := Parse(s); err == nil {
			t.Errorf("Parse(%q): want error", s)
		}
	}
}
