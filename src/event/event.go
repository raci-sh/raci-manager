// Package event defines the append-only audit events of SPEC.md section 4.1:
// a common envelope, a type discriminator and one payload per type.
package event

import (
	"encoding/json"
	"fmt"

	"github.com/raci-sh/raci-manager/src/canonical"
	"github.com/raci-sh/raci-manager/src/digest"
	"github.com/raci-sh/raci-manager/src/model"
)

// Type selects the payload of an event.
type Type string

// The event types.
const (
	TypeEntry      Type = "entry"
	TypeResource   Type = "resource"
	TypeConstraint Type = "constraint"
	TypeRule       Type = "rule"
	TypePin        Type = "pin"
	TypeCommit     Type = "commit"
)

// Via says how a change arrived.
type Via string

// The ways a change can arrive.
const (
	ViaUser      Via = "user"
	ViaAgent     Via = "agent"
	ViaImport    Via = "import"
	ViaInference Via = "inference"
)

// Actor is the party behind an event and how its change arrived.
type Actor struct {
	Party model.PartyID `json:"party"`
	Via   Via           `json:"via"`
}

// Cause links an inferred event to the rule and the original event that
// produced it. Inferred events are never causes.
type Cause struct {
	RuleID  string `json:"rule_id"`
	EventID string `json:"event_id"`
}

// Status is the outcome of a change.
type Status string

// The outcomes.
const (
	Applied  Status = "applied"
	Rejected Status = "rejected"
)

// Outcome is applied or rejected, with a reason when rejected.
type Outcome struct {
	Status Status `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// ConstraintResult records one enabled constraint evaluated for the event.
type ConstraintResult struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
}

// Event is one line of a realm's log. Payload holds the type-specific body
// (EntryPayload, ResourcePayload, ...). The exact JSON layout of the payload
// within the envelope is a scaffold choice, not yet fixed by the spec.
type Event struct {
	ID          string             `json:"id"`
	Seq         uint64             `json:"seq"`
	TS          string             `json:"ts"`
	Realm       model.RealmID      `json:"realm"`
	Type        Type               `json:"type"`
	ChangeID    string             `json:"change_id"`
	Actor       Actor              `json:"actor"`
	Cause       *Cause             `json:"cause,omitempty"`
	Constraints []ConstraintResult `json:"constraints,omitempty"`
	Outcome     Outcome            `json:"outcome"`
	Payload     json.RawMessage    `json:"payload"`
	PrevHash    digest.Digest      `json:"prev_hash"`
	Hash        digest.Digest      `json:"hash,omitempty"`
}

// ComputeHash returns blake3 over the canonical JSON of the event without
// its hash field and with prev_hash included.
func (e *Event) ComputeHash() (digest.Digest, error) {
	c := *e
	c.Hash = ""
	b, err := canonical.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("event %s: %w", e.ID, err)
	}
	return digest.Sum(b), nil
}

// Seal links the event to prev and sets its hash.
func (e *Event) Seal(prev digest.Digest) error {
	e.PrevHash = prev
	h, err := e.ComputeHash()
	if err != nil {
		return err
	}
	e.Hash = h
	return nil
}

// SetPayload encodes p as the event payload.
func (e *Event) SetPayload(p any) error {
	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("event %s: payload: %w", e.ID, err)
	}
	e.Payload = b
	return nil
}

// EntryOp is the operation of an entry event.
type EntryOp string

// The entry operations.
const (
	EntryAdd     EntryOp = "add"
	EntryUpdate  EntryOp = "update"
	EntryRemove  EntryOp = "remove"
	EntryReplace EntryOp = "replace"
)

// Target names the entry an entry event changes.
type Target struct {
	Activity model.ActivityID `json:"activity"`
	Resource model.ResourceID `json:"resource"`
	Party    model.PartyID    `json:"party"`
}

// EntryPayload is the payload of a TypeEntry event.
type EntryPayload struct {
	Op       EntryOp     `json:"op"`
	Target   Target      `json:"target"`
	Before   model.Roles `json:"before"`
	After    model.Roles `json:"after"`
	Conflict []string    `json:"conflict,omitempty"`
}

// ResourceOp is the operation of a resource event.
type ResourceOp string

// The resource operations.
const (
	ResourceNew    ResourceOp = "new"
	ResourceDelete ResourceOp = "delete"
)

// ResourcePayload is the payload of a TypeResource event.
type ResourcePayload struct {
	Op       ResourceOp     `json:"op"`
	Resource model.Resource `json:"resource"`
	Realm    model.RealmID  `json:"realm"`
}

// ConstraintOp is the operation of a constraint event. There is no disable
// or remove: constraints are tighten-only.
type ConstraintOp string

// The constraint operations.
const (
	ConstraintEnable  ConstraintOp = "enable"
	ConstraintAddPair ConstraintOp = "add-pair"
)

// ConstraintPayload is the payload of a TypeConstraint event.
type ConstraintPayload struct {
	Constraint string          `json:"constraint"`
	Op         ConstraintOp    `json:"op"`
	Params     json.RawMessage `json:"params,omitempty"`
}

// Transition is a SteerSpec lifecycle state change (D, A, P, I, R, T).
type Transition struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// RulePayload is the payload of a TypeRule event.
type RulePayload struct {
	RuleID     string          `json:"rule_id"`
	Revision   int             `json:"revision"`
	Transition *Transition     `json:"transition,omitempty"`
	Supersedes string          `json:"supersedes,omitempty"`
	Body       json.RawMessage `json:"body"`
	BodyHash   digest.Digest   `json:"body_hash"`
}

// PinPayload is the payload of a TypePin event.
type PinPayload struct {
	DependencyRealm model.RealmID `json:"dependency_realm"`
	Version         string        `json:"version"`
	Op              string        `json:"op"`
}

// Finding is a non-blocking finding recorded on a commit marker.
type Finding struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// CommitPayload is the payload of a TypeCommit event, appended last for each
// change. Replay applies a change only when its marker is applied.
type CommitPayload struct {
	ChangeID string    `json:"change_id"`
	Events   []string  `json:"events"`
	Outcome  Outcome   `json:"outcome"`
	Findings []Finding `json:"findings,omitempty"`
}
