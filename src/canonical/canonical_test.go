package canonical

import "testing"

func TestMarshalSortsKeysAndDropsWhitespace(t *testing.T) {
	in := map[string]any{"realm": "sh.raci", "raci": "1", "log": map[string]any{"head_seq": 0, "head_hash": "x"}}
	got, err := Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"log":{"head_hash":"x","head_seq":0},"raci":"1","realm":"sh.raci"}`
	if string(got) != want {
		t.Errorf("got %s\nwant %s", got, want)
	}
}

func TestTransformRejectsInvalidJSON(t *testing.T) {
	if _, err := Transform([]byte("{")); err == nil {
		t.Error("want error for invalid JSON")
	}
	if _, err := Marshal(make(chan int)); err == nil {
		t.Error("want error for unmarshalable value")
	}
}
