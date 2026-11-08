// Package inference defines the inference rule vocabulary of SPEC.md
// section 3.1, taken from ESIS: a rule fires on an operation, for a target
// resource type, and produces role deltas through map, add-issuer and strip
// actions. Inference is non-recursive: inferred changes never fire rules.
package inference

import (
	"errors"
	"fmt"

	"github.com/raci-sh/raci-manager/src/model"
)

// Operation is the change that fires a rule.
type Operation string

// Resource events (new, delete) and entry events (add, update, remove).
const (
	OpNew    Operation = "new"
	OpDelete Operation = "delete"
	OpAdd    Operation = "add"
	OpUpdate Operation = "update"
	OpRemove Operation = "remove"
)

// IsEntryOp reports whether op is an entry event operation.
func (op Operation) IsEntryOp() bool { return op == OpAdd || op == OpUpdate || op == OpRemove }

// Valid reports whether op is a known operation.
func (op Operation) Valid() bool { return op == OpNew || op == OpDelete || op.IsEntryOp() }

// ActionKind is one of the action vocabulary words.
type ActionKind string

// The actions.
const (
	ActionMap       ActionKind = "map"
	ActionAddIssuer ActionKind = "add-issuer"
	ActionStrip     ActionKind = "strip"
)

// Action is one action of a rule. Which fields apply depends on Kind:
//   - map(SourceType, FromRole, ToRole)
//   - add-issuer(Role)
//   - strip(SourceType, HeldRole, Role)
type Action struct {
	Kind       ActionKind `json:"kind"`
	SourceType string     `json:"source_type,omitempty"`
	FromRole   model.Role `json:"from_role,omitempty"`
	ToRole     model.Role `json:"to_role,omitempty"`
	HeldRole   model.Role `json:"held_role,omitempty"`
	Role       model.Role `json:"role,omitempty"`
	// AddOnly adds missing entries and never changes an existing one.
	AddOnly bool `json:"add_only,omitempty"`
}

// Rule is an inference rule. Rules are evaluated in rule-id order.
type Rule struct {
	ID          string    `json:"id"`
	TargetType  string    `json:"target_type"`
	Operation   Operation `json:"operation"`
	TriggerType string    `json:"trigger_type,omitempty"`
	// Selector names the source resource explicitly. Syntax is open.
	Selector string   `json:"selector,omitempty"`
	Actions  []Action `json:"actions"`
}

// Validate checks the structural rules of SPEC.md section 3.1.
func (r Rule) Validate() error {
	var errs []error
	if r.ID == "" || r.TargetType == "" {
		errs = append(errs, errors.New("id and target type are required"))
	}
	if !r.Operation.Valid() {
		errs = append(errs, fmt.Errorf("unknown operation %q", r.Operation))
	}
	if r.TriggerType != "" && !r.Operation.IsEntryOp() {
		errs = append(errs, fmt.Errorf("a rule with a trigger type needs an entry operation, not %q", r.Operation))
	}
	if len(r.Actions) == 0 {
		errs = append(errs, errors.New("at least one action is required"))
	}
	for i, a := range r.Actions {
		if err := a.validate(r.Operation); err != nil {
			errs = append(errs, fmt.Errorf("action %d: %w", i, err))
		}
	}
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("rule %q: %w", r.ID, err)
	}
	return nil
}

func (a Action) validate(op Operation) error {
	switch a.Kind {
	case ActionMap:
		if a.SourceType == "" || !a.FromRole.Valid() || !a.ToRole.Valid() {
			return errors.New("map needs source type, from role and to role")
		}
	case ActionAddIssuer:
		if !a.Role.Valid() {
			return errors.New("add-issuer needs a role")
		}
		if op != OpNew {
			return errors.New("add-issuer is only valid in rules whose operation is new")
		}
	case ActionStrip:
		if a.SourceType == "" || !a.HeldRole.Valid() || !a.Role.Valid() {
			return errors.New("strip needs source type, held role and role")
		}
	default:
		return fmt.Errorf("unknown action %q", a.Kind)
	}
	return nil
}
