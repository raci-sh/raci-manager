package eventlog

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/raci-sh/raci-manager/src/digest"
	"github.com/raci-sh/raci-manager/src/event"
)

func newEvent(id string) *event.Event {
	return &event.Event{ID: id, Realm: "sh.raci", Type: event.TypeCommit, ChangeID: "c1",
		Actor: event.Actor{Party: "alice", Via: event.ViaUser}, Outcome: event.Outcome{Status: event.Applied},
		Payload: []byte(`{}`)}
}

func TestAppendAndVerify(t *testing.T) {
	log := Open(filepath.Join(t.TempDir(), FileName))
	events, err := log.Read()
	if err != nil || len(events) != 0 || HeadOf(events) != EmptyHead {
		t.Fatalf("empty log: %v %v", events, err)
	}
	for _, id := range []string{"e1", "e2"} {
		if err := log.Append(newEvent(id)); err != nil {
			t.Fatal(err)
		}
	}
	events, err = log.Read()
	if err != nil || len(events) != 2 {
		t.Fatalf("read: %v %v", events, err)
	}
	if events[0].PrevHash != digest.Zero || events[1].PrevHash != events[0].Hash || events[1].Seq != 2 {
		t.Fatalf("chain not linked: %+v", events)
	}
	if err := VerifyChain(events); err != nil {
		t.Fatal(err)
	}
	if HeadOf(events).Seq != 2 || log.Path() == "" {
		t.Error("bad head or path")
	}
}

func TestVerifyChainDetectsTampering(t *testing.T) {
	log := Open(filepath.Join(t.TempDir(), FileName))
	for _, id := range []string{"e1", "e2", "e3"} {
		if err := log.Append(newEvent(id)); err != nil {
			t.Fatal(err)
		}
	}
	events, _ := log.Read()

	tampered := append([]event.Event(nil), events...)
	tampered[1].ChangeID = "forged"
	var ce *ChainError
	if err := VerifyChain(tampered); !errors.As(err, &ce) || ce.Seq != 2 || !strings.Contains(err.Error(), "seq 2") {
		t.Errorf("content tamper: got %v", err)
	}

	relinked := append([]event.Event(nil), events...)
	relinked[2].PrevHash = digest.Zero
	if err := VerifyChain(relinked); !errors.As(err, &ce) || ce.Seq != 3 {
		t.Errorf("link tamper: got %v", err)
	}

	if err := VerifyChain(events[1:]); !errors.As(err, &ce) || ce.Seq != 2 {
		t.Errorf("gap: got %v", err)
	}
}

func TestReadRejectsGarbage(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	if err := os.WriteFile(path, []byte("{\"id\":\"e1\"}\n\nnot json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(path).Read(); err == nil || !strings.Contains(err.Error(), "line 3") {
		t.Errorf("want line 3 error, got %v", err)
	}
	if err := Open(filepath.Join(path, "sub")).Append(newEvent("x")); err == nil {
		t.Error("want error appending under a file")
	}
}
