package inference

import (
	"strings"
	"testing"

	"github.com/raci-sh/raci-manager/src/model"
)

func TestValidateESISExample(t *testing.T) {
	r := Rule{ID: "INF-001", TargetType: "audit", Operation: OpNew, Actions: []Action{
		{Kind: ActionMap, SourceType: "module", FromRole: model.Responsible, ToRole: model.Accountable},
		{Kind: ActionMap, SourceType: "module", FromRole: model.Accountable, ToRole: model.Informed},
		{Kind: ActionAddIssuer, Role: model.Responsible},
		{Kind: ActionStrip, SourceType: "module", HeldRole: model.Informed, Role: model.Accountable},
	}}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateRejects(t *testing.T) {
	cases := map[string]Rule{
		"trigger type needs an entry operation": {ID: "x", TargetType: "audit", Operation: OpNew, TriggerType: "module",
			Actions: []Action{{Kind: ActionAddIssuer, Role: model.Responsible}}},
		"add-issuer is only valid": {ID: "x", TargetType: "audit", Operation: OpAdd,
			Actions: []Action{{Kind: ActionAddIssuer, Role: model.Responsible}}},
		"unknown operation":       {ID: "x", TargetType: "audit", Operation: "touch", Actions: []Action{{Kind: "copy"}}},
		"at least one action":     {ID: "x", TargetType: "audit", Operation: OpDelete},
		"id and target type":      {Operation: OpRemove, Actions: []Action{{Kind: ActionMap}}},
		"strip needs source type": {ID: "x", TargetType: "audit", Operation: OpUpdate, Actions: []Action{{Kind: ActionStrip}}},
		"add-issuer needs a role": {ID: "x", TargetType: "audit", Operation: OpNew, Actions: []Action{{Kind: ActionAddIssuer}}},
		"unknown action":          {ID: "x", TargetType: "audit", Operation: OpNew, Actions: []Action{{Kind: "copy"}}},
	}
	for want, r := range cases {
		err := r.Validate()
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("want %q, got %v", want, err)
		}
	}
}
