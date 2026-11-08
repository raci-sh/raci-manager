package view

import (
	"errors"
	"strings"
	"testing"

	"github.com/raci-sh/raci-manager/src/constraint"
	"github.com/raci-sh/raci-manager/src/event"
	"github.com/raci-sh/raci-manager/src/model"
)

func TestEmptyViewBytes(t *testing.T) {
	v := Empty("sh.raci")
	b, err := v.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := `{"constraints":{"enabled":[],"pairs":[]},"entries":[],"log":{"head_hash":"blake3:` +
		strings.Repeat("0", 64) + `","head_seq":0},"pins":[],"raci":"1","realm":"sh.raci","resources":[],"rules_in_force":[]}`
	if string(b) != want {
		t.Errorf("got  %s\nwant %s", b, want)
	}
}

func TestNormalizeOrdersArrays(t *testing.T) {
	v := Empty("sh.raci")
	v.Entries = []model.Entry{
		{Activity: "update", Party: "bob", Resource: "billing", Roles: "R"},
		{Activity: "add", Party: "zed", Resource: "billing", Roles: "I"},
		{Activity: "add", Party: "amy", Resource: "billing", Roles: "A"},
		{Activity: "add", Party: "amy", Resource: "audit", Roles: "C"},
	}
	v.Resources = []model.Resource{{ID: "billing", Type: "module"}, {ID: "audit", Type: "audit", Parent: "billing"}}
	v.Pins = []Pin{{DependencyRealm: "sh.raci", Version: "2"}, {DependencyRealm: "dev.steerspec", Version: "1"}, {DependencyRealm: "sh.raci", Version: "1"}}
	v.RulesInForce = []RuleInForce{{RuleID: "RACIINF-002"}, {RuleID: "RACIINF-001"}}
	v.Constraints.Enabled = []constraint.Name{constraint.SingleResponsible, constraint.MinRA}
	v.Constraints.Pairs = []constraint.Pair{
		{A: constraint.Side{Activity: "b", Roles: "R"}},
		{A: constraint.Side{Activity: "a", Roles: "R"}},
	}
	if err := v.Normalize(); err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range v.Entries {
		got = append(got, string(e.Resource)+"/"+string(e.Activity)+"/"+string(e.Party))
	}
	if strings.Join(got, ",") != "audit/add/amy,billing/add/amy,billing/add/zed,billing/update/bob" {
		t.Errorf("entries: %v", got)
	}
	if v.Resources[0].ID != "audit" || v.Pins[0].DependencyRealm != "dev.steerspec" || v.Pins[2].Version != "2" ||
		v.RulesInForce[0].RuleID != "RACIINF-001" || v.Constraints.Enabled[0] != constraint.MinRA ||
		v.Constraints.Pairs[0].A.Activity != "a" {
		t.Errorf("not normalized: %+v", v)
	}
}

func TestReplay(t *testing.T) {
	v, err := Replay("sh.raci", nil)
	if err != nil || v.Realm != "sh.raci" {
		t.Errorf("empty replay: %+v %v", v, err)
	}
	if _, err := Replay("sh.raci", []event.Event{{ID: "e1"}}); !errors.Is(err, ErrNotImplemented) {
		t.Errorf("want ErrNotImplemented, got %v", err)
	}
}
