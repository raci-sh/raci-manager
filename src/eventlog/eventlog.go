// Package eventlog stores a realm's audit log offline: one JSONL file per
// realm, next to its .raci file, with the CLI as its only writer (SPEC.md
// section 4.2). Each line is one event in canonical JSON, and the events form
// a hash chain through prev_hash.
package eventlog

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/raci-sh/raci-manager/src/canonical"
	"github.com/raci-sh/raci-manager/src/digest"
	"github.com/raci-sh/raci-manager/src/event"
)

// FileName is the conventional log name next to a .raci file. The spec
// leaves the name open; this is the scaffold's proposal.
const FileName = ".raci.jsonl"

// Log is a per-realm JSONL event log.
type Log struct {
	path string
}

// Open returns the log at path. The file need not exist yet: a missing log is
// an empty log.
func Open(path string) *Log { return &Log{path: path} }

// Path returns the file path of the log.
func (l *Log) Path() string { return l.path }

// Read returns every event in the log, in file order.
func (l *Log) Read() ([]event.Event, error) {
	data, err := os.ReadFile(l.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("eventlog: %w", err)
	}
	var events []event.Event
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for line := 1; sc.Scan(); line++ {
		if len(bytes.TrimSpace(sc.Bytes())) == 0 {
			continue
		}
		var e event.Event
		if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
			return nil, fmt.Errorf("eventlog: %s line %d: %w", l.path, line, err)
		}
		events = append(events, e)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("eventlog: %w", err)
	}
	return events, nil
}

// Head is the position of the last event of a log.
type Head struct {
	Seq  uint64        `json:"seq"`
	Hash digest.Digest `json:"hash"`
}

// EmptyHead is the head of an empty log: seq 0 and the all-zero digest.
var EmptyHead = Head{Seq: 0, Hash: digest.Zero}

// HeadOf returns the head of events, or EmptyHead for none.
func HeadOf(events []event.Event) Head {
	if len(events) == 0 {
		return EmptyHead
	}
	last := events[len(events)-1]
	return Head{Seq: last.Seq, Hash: last.Hash}
}

// Append sets e's seq and prev_hash from the current head, seals it and
// writes it as one canonical JSON line. Offline mode is single-writer per
// realm; Append does not lock.
func (l *Log) Append(e *event.Event) error {
	events, err := l.Read()
	if err != nil {
		return err
	}
	head := HeadOf(events)
	e.Seq = head.Seq + 1
	if err := e.Seal(head.Hash); err != nil {
		return err
	}
	line, err := canonical.Marshal(e)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("eventlog: %w", err)
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		_ = f.Close()
		return fmt.Errorf("eventlog: %w", err)
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return fmt.Errorf("eventlog: %w", err)
	}
	return f.Close()
}

// ChainError reports the first event at which the hash chain breaks.
type ChainError struct {
	Seq    uint64
	Reason string
}

func (e *ChainError) Error() string {
	return fmt.Sprintf("hash chain broken at seq %d: %s", e.Seq, e.Reason)
}

// VerifyChain recomputes every hash and checks every prev_hash link and seq,
// failing at the first mismatch with a *ChainError.
func VerifyChain(events []event.Event) error {
	prev := EmptyHead
	for i := range events {
		e := &events[i]
		if e.Seq != prev.Seq+1 {
			return &ChainError{Seq: e.Seq, Reason: fmt.Sprintf("seq %d follows %d", e.Seq, prev.Seq)}
		}
		if e.PrevHash != prev.Hash {
			return &ChainError{Seq: e.Seq, Reason: "prev_hash does not match the previous event's hash"}
		}
		got, err := e.ComputeHash()
		if err != nil {
			return &ChainError{Seq: e.Seq, Reason: err.Error()}
		}
		if got != e.Hash {
			return &ChainError{Seq: e.Seq, Reason: "hash does not match the event content"}
		}
		prev = Head{Seq: e.Seq, Hash: e.Hash}
	}
	return nil
}
