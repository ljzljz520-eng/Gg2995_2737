package audit

import (
	"lawdrive/internal/model"
	"testing"
)

func TestAuditSummaries(t *testing.T) {
	entries := []model.AuditEntry{{ActorID: "u1", Action: "view", OccurredAt: 10}, {ActorID: "u1", Action: "edit", OccurredAt: 20}, {ActorID: "u2", Action: "view", OccurredAt: 15}}
	reports := Summarize(entries)
	if len(reports) != 2 || reports[0].Actions != 2 {
		t.Fatalf("reports=%v", reports)
	}
}
