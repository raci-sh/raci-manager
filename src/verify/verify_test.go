package verify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/raci-sh/raci-manager/src/event"
	"github.com/raci-sh/raci-manager/src/eventlog"
)

func TestEverySeverityIsSet(t *testing.T) {
	for _, c := range []Code{HashMismatch, ViewMismatch, StaleView, UnrecordedRule, IncompleteChange, NoMatchingRule, DescendantViolation} {
		if DefaultSeverity[c] == "" {
			t.Errorf("no severity for %q", c)
		}
	}
}

func TestLog(t *testing.T) {
	dir := t.TempDir()
	log := eventlog.Open(filepath.Join(dir, eventlog.FileName))
	r, err := Log("sh.raci", log)
	if err != nil || r.HasErrors() || r.LogHead != eventlog.EmptyHead {
		t.Fatalf("empty log: %+v %v", r, err)
	}
	if err := log.Append(&event.Event{ID: "e1", Type: event.TypeCommit, Payload: []byte(`{}`)}); err != nil {
		t.Fatal(err)
	}
	r, _ = Log("sh.raci", log)
	if r.HasErrors() || r.LogHead.Seq != 1 {
		t.Fatalf("one event: %+v", r)
	}
	data, _ := os.ReadFile(log.Path())
	if err := os.WriteFile(log.Path(), append([]byte(`{"seq":1,"id":"forged","prev_hash":"x"}`+"\n"), data...), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err = Log("sh.raci", log)
	if err != nil || !r.HasErrors() || r.Findings[0].Code != HashMismatch {
		t.Errorf("tampered: %+v %v", r, err)
	}
	if err := os.WriteFile(log.Path(), []byte("nope\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Log("sh.raci", log); err == nil {
		t.Error("want read error")
	}
}
