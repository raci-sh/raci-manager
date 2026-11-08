package model

import "testing"

func TestParseRoles(t *testing.T) {
	cases := map[string]Roles{"": "", "RA": "AR", "ARA": "AR", "ICRA": "ARCI", "c": ""}
	for in, want := range cases {
		got, err := ParseRoles(in)
		if in == "c" {
			if err == nil {
				t.Errorf("ParseRoles(%q): want error", in)
			}
			continue
		}
		if err != nil || got != want {
			t.Errorf("ParseRoles(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}

func TestRolesHasAndBest(t *testing.T) {
	rs := Roles("RC")
	if !rs.Has(Responsible) || rs.Has(Accountable) {
		t.Errorf("Has: unexpected result for %q", rs)
	}
	if b, ok := rs.Best(); !ok || b != Responsible || b.String() != "R" {
		t.Errorf("Best(%q) = %v, %v", rs, b, ok)
	}
	if _, ok := Roles("").Best(); ok {
		t.Error("Best of empty set should be false")
	}
}
