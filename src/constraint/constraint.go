// Package constraint names the strictness constraints of SPEC.md section 2.
// Constraints are switched on per realm and are tighten-only: a sub-realm may
// add constraints or conflict pairs but never remove an inherited one.
package constraint

import (
	"slices"

	"github.com/raci-sh/raci-manager/src/model"
)

// Name is the name of a constraint.
type Name string

// The constraints of SPEC.md section 2.
const (
	SingleAccountable  Name = "single-accountable"
	SingleResponsible  Name = "single-responsible"
	OneRolePerParty    Name = "one-role-per-party"
	MinRA              Name = "min-ra"
	SeparationOfDuties Name = "separation-of-duties"
)

// All lists every known constraint, alphabetically.
var All = []Name{MinRA, OneRolePerParty, SeparationOfDuties, SingleAccountable, SingleResponsible}

// Strict is the "strict" preset. separation-of-duties is opt-in because it
// needs a conflict list.
var Strict = []Name{MinRA, OneRolePerParty, SingleAccountable, SingleResponsible}

// Known reports whether n is a known constraint name.
func Known(n Name) bool { return slices.Contains(All, n) }

// Side is one half of a separation-of-duties conflict pair. Selector syntax
// for activities and resources is open (SPEC.md section 2.1).
type Side struct {
	Activity string      `json:"activity"`
	Resource string      `json:"resource"`
	Roles    model.Roles `json:"roles"`
}

// Pair declares that a party holding A.Roles on (A.Activity, A.Resource) must
// not hold B.Roles on (B.Activity, B.Resource).
type Pair struct {
	A Side `json:"a"`
	B Side `json:"b"`
}
