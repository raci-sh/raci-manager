package event

import (
	"testing"

	"github.com/raci-sh/raci-manager/src/digest"
)

func TestSealIsDeterministicAndChained(t *testing.T) {
	e := Event{ID: "e1", Seq: 1, Realm: "sh.raci", Type: TypeEntry, ChangeID: "c1",
		Actor: Actor{Party: "alice", Via: ViaUser}, Outcome: Outcome{Status: Applied}}
	if err := e.SetPayload(EntryPayload{Op: EntryAdd, Target: Target{Activity: "update", Resource: "billing", Party: "bob"}, After: "R"}); err != nil {
		t.Fatal(err)
	}
	if err := e.Seal(digest.Zero); err != nil {
		t.Fatal(err)
	}
	if !e.Hash.Valid() || e.PrevHash != digest.Zero {
		t.Fatalf("bad seal: %+v", e)
	}
	again, err := e.ComputeHash()
	if err != nil || again != e.Hash {
		t.Errorf("hash not reproducible: %s vs %s (%v)", again, e.Hash, err)
	}
	e.Seq = 2
	if changed, _ := e.ComputeHash(); changed == e.Hash {
		t.Error("hash did not change with the content")
	}
}

func TestSetPayloadError(t *testing.T) {
	var e Event
	if err := e.SetPayload(make(chan int)); err == nil {
		t.Error("want error")
	}
}
