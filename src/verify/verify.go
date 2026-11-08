// Package verify holds the verification report of SPEC.md section 7.1.
// Findings are computed on demand and never persisted in the .raci view.
package verify

import (
	"github.com/raci-sh/raci-manager/src/eventlog"
	"github.com/raci-sh/raci-manager/src/model"
)

// Code identifies a kind of finding.
type Code string

// The finding codes.
const (
	HashMismatch        Code = "hash-mismatch"
	ViewMismatch        Code = "view-mismatch"
	StaleView           Code = "stale-view"
	UnrecordedRule      Code = "unrecorded-rule"
	IncompleteChange    Code = "incomplete-change"
	NoMatchingRule      Code = "no-matching-rule"
	DescendantViolation Code = "descendant-violation"
)

// Severity of a finding.
type Severity string

// The severities.
const (
	Error   Severity = "error"
	Warning Severity = "warning"
)

// DefaultSeverity maps each code to its severity.
var DefaultSeverity = map[Code]Severity{
	HashMismatch:        Error,
	ViewMismatch:        Error,
	StaleView:           Warning,
	UnrecordedRule:      Error,
	IncompleteChange:    Warning,
	NoMatchingRule:      Warning,
	DescendantViolation: Error,
}

// Finding is one verification finding.
type Finding struct {
	Code     Code          `json:"code"`
	Severity Severity      `json:"severity"`
	Realm    model.RealmID `json:"realm"`
	Detail   string        `json:"detail"`
}

// Report is the output of `raci verify`.
type Report struct {
	Realm    model.RealmID `json:"realm"`
	LogHead  eventlog.Head `json:"log_head"`
	Findings []Finding     `json:"findings"`
}

// NewReport returns an empty report for realm at head.
func NewReport(realm model.RealmID, head eventlog.Head) *Report {
	return &Report{Realm: realm, LogHead: head, Findings: []Finding{}}
}

// Add records a finding with the code's default severity.
func (r *Report) Add(code Code, realm model.RealmID, detail string) {
	r.Findings = append(r.Findings, Finding{Code: code, Severity: DefaultSeverity[code], Realm: realm, Detail: detail})
}

// HasErrors reports whether any finding is an error.
func (r *Report) HasErrors() bool {
	for _, f := range r.Findings {
		if f.Severity == Error {
			return true
		}
	}
	return false
}

// Log verifies the hash chain of a realm's log and returns the report.
// View, rule and change checks are not implemented yet.
func Log(realm model.RealmID, log *eventlog.Log) (*Report, error) {
	events, err := log.Read()
	if err != nil {
		return nil, err
	}
	r := NewReport(realm, eventlog.HeadOf(events))
	if err := eventlog.VerifyChain(events); err != nil {
		r.Add(HashMismatch, realm, err.Error())
	}
	return r, nil
}
