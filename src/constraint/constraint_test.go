package constraint

import "testing"

func TestStrictIsKnownAndExcludesSoD(t *testing.T) {
	for _, n := range Strict {
		if !Known(n) {
			t.Errorf("%q in Strict but not known", n)
		}
		if n == SeparationOfDuties {
			t.Error("strict must not enable separation-of-duties")
		}
	}
	if Known("auto-demote") {
		t.Error("unknown name reported as known")
	}
}
