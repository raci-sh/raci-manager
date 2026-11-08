// Package model holds the RACI core model (SPEC.md section 1): realms,
// parties, activities, resources and the entries that bind a party to a set
// of roles for an activity on a resource.
package model

import (
	"fmt"
	"strings"
)

// Role is one of the four RACI roles.
type Role byte

// The four roles, in precedence order (A > R > C > I).
const (
	Accountable Role = 'A'
	Responsible Role = 'R'
	Consulted   Role = 'C'
	Informed    Role = 'I'
)

// Precedence lists the roles from strongest to weakest. It is also the
// order in which roles are written in a Roles string.
var Precedence = []Role{Accountable, Responsible, Consulted, Informed}

// String returns the single-letter form of the role.
func (r Role) String() string { return string(r) }

// Valid reports whether r is one of A, R, C or I.
func (r Role) Valid() bool {
	return r == Accountable || r == Responsible || r == Consulted || r == Informed
}

// Roles is a set of roles written as a string in precedence order, such as
// "AR". The empty string is the empty set. The data model is permissive: a
// party may hold several roles on the same activity and resource; strictness
// is a constraint (package constraint), not part of the schema.
type Roles string

// ParseRoles validates s and returns it in canonical order without
// duplicates, so "RA" and "ARA" both become "AR".
func ParseRoles(s string) (Roles, error) {
	var seen [4]bool
	for i := 0; i < len(s); i++ {
		r := Role(s[i])
		if !r.Valid() {
			return "", fmt.Errorf("invalid role %q in %q", s[i], s)
		}
		seen[index(r)] = true
	}
	var b strings.Builder
	for i, r := range Precedence {
		if seen[i] {
			b.WriteByte(byte(r))
		}
	}
	return Roles(b.String()), nil
}

// Has reports whether the set contains role.
func (rs Roles) Has(role Role) bool {
	return strings.IndexByte(string(rs), byte(role)) >= 0
}

// Best returns the strongest role in the set, and false for an empty set.
func (rs Roles) Best() (Role, bool) {
	if rs == "" {
		return 0, false
	}
	return Role(rs[0]), true
}

func index(r Role) int {
	for i, p := range Precedence {
		if p == r {
			return i
		}
	}
	return -1
}

// RealmID identifies a realm, matching a SteerSpec realm id (e.g. "sh.raci").
type RealmID string

// PartyID identifies a human or agent party.
type PartyID string

// ActivityID names an activity carried out on a resource (e.g. "update").
type ActivityID string

// ResourceID identifies a resource within a realm.
type ResourceID string

// Resource is something activities are carried out on. Its realm is carried
// by the event or view that holds it.
type Resource struct {
	ID     ResourceID `json:"id"`
	Parent ResourceID `json:"parent,omitempty"`
	Type   string     `json:"type"`
}

// Entry binds a party to a set of roles for an activity on a resource.
type Entry struct {
	Activity ActivityID `json:"activity"`
	Party    PartyID    `json:"party"`
	Resource ResourceID `json:"resource"`
	Roles    Roles      `json:"roles"`
}
