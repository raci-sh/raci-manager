// Package view is the .raci file (SPEC.md section 7): a materialised view of
// one realm's audit log, in canonical JSON, rebuilt and verified by replay.
// It holds RACI and nothing else: no findings, no history, no paths.
package view

import (
	"bytes"
	"errors"
	"slices"
	"strings"

	"github.com/raci-sh/raci-manager/src/canonical"
	"github.com/raci-sh/raci-manager/src/constraint"
	"github.com/raci-sh/raci-manager/src/digest"
	"github.com/raci-sh/raci-manager/src/event"
	"github.com/raci-sh/raci-manager/src/model"
)

// FileName is the conventional name of the view file.
const FileName = ".raci"

// FormatVersion is the value of the "raci" member.
const FormatVersion = "1"

// ErrNotImplemented is returned by parts of the view not built yet.
var ErrNotImplemented = errors.New("view: not implemented yet")

// Constraints are the constraints recorded by this realm's own events. The
// effective set along the realm path is computed at load, never embedded.
type Constraints struct {
	Enabled []constraint.Name `json:"enabled"`
	Pairs   []constraint.Pair `json:"pairs"`
}

// LogHead identifies the log state the view reflects.
type LogHead struct {
	HeadHash digest.Digest `json:"head_hash"`
	HeadSeq  uint64        `json:"head_seq"`
}

// Pin is the version of a dependency realm this realm is evaluated against.
type Pin struct {
	DependencyRealm model.RealmID `json:"dependency_realm"`
	Version         string        `json:"version"`
}

// RuleInForce is a rule recorded in force in this realm's log.
type RuleInForce struct {
	BodyHash digest.Digest `json:"body_hash"`
	Revision int           `json:"revision"`
	RuleID   string        `json:"rule_id"`
}

// View is the content of a .raci file.
type View struct {
	Constraints  Constraints      `json:"constraints"`
	Entries      []model.Entry    `json:"entries"`
	Log          LogHead          `json:"log"`
	Pins         []Pin            `json:"pins"`
	Raci         string           `json:"raci"`
	Realm        model.RealmID    `json:"realm"`
	Resources    []model.Resource `json:"resources"`
	RulesInForce []RuleInForce    `json:"rules_in_force"`
}

// Empty returns the view of an empty log for realm.
func Empty(realm model.RealmID) View {
	return View{
		Constraints:  Constraints{Enabled: []constraint.Name{}, Pairs: []constraint.Pair{}},
		Entries:      []model.Entry{},
		Log:          LogHead{HeadHash: digest.Zero, HeadSeq: 0},
		Pins:         []Pin{},
		Raci:         FormatVersion,
		Realm:        realm,
		Resources:    []model.Resource{},
		RulesInForce: []RuleInForce{},
	}
}

// Normalize sorts every array into its canonical order (SPEC.md section 7).
func (v *View) Normalize() error {
	slices.SortFunc(v.Resources, func(a, b model.Resource) int { return strings.Compare(string(a.ID), string(b.ID)) })
	slices.SortFunc(v.Entries, func(a, b model.Entry) int {
		if c := strings.Compare(string(a.Resource), string(b.Resource)); c != 0 {
			return c
		}
		if c := strings.Compare(string(a.Activity), string(b.Activity)); c != 0 {
			return c
		}
		return strings.Compare(string(a.Party), string(b.Party))
	})
	slices.SortFunc(v.RulesInForce, func(a, b RuleInForce) int { return strings.Compare(a.RuleID, b.RuleID) })
	slices.SortFunc(v.Pins, func(a, b Pin) int {
		if c := strings.Compare(string(a.DependencyRealm), string(b.DependencyRealm)); c != 0 {
			return c
		}
		return strings.Compare(a.Version, b.Version)
	})
	slices.Sort(v.Constraints.Enabled)
	type keyed struct {
		p constraint.Pair
		k []byte
	}
	ks := make([]keyed, len(v.Constraints.Pairs))
	for i, p := range v.Constraints.Pairs {
		k, err := canonical.Marshal(p)
		if err != nil {
			return err
		}
		ks[i] = keyed{p, k}
	}
	slices.SortFunc(ks, func(a, b keyed) int { return bytes.Compare(a.k, b.k) })
	for i := range ks {
		v.Constraints.Pairs[i] = ks[i].p
	}
	return nil
}

// Marshal normalizes v and returns its canonical JSON, the exact bytes of a
// .raci file. The same log always yields byte-identical files.
func (v *View) Marshal() ([]byte, error) {
	if err := v.Normalize(); err != nil {
		return nil, err
	}
	return canonical.Marshal(v)
}

// Replay rebuilds the view of realm from its log. Not implemented yet: it
// will apply only changes whose commit marker is present and applied.
func Replay(realm model.RealmID, events []event.Event) (View, error) {
	if len(events) == 0 {
		return Empty(realm), nil
	}
	return View{}, ErrNotImplemented
}
