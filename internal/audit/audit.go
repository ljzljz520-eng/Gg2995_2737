package audit

import (
	"fmt"
	"lawdrive/internal/model"
)

func New(id, caseID, actorID, action, targetID string, detail map[string]string, at int64) (model.AuditEntry, error) {
	if id == "" || caseID == "" || actorID == "" || action == "" || targetID == "" {
		return model.AuditEntry{}, fmt.Errorf("audit identity is incomplete")
	}
	if at <= 0 {
		return model.AuditEntry{}, fmt.Errorf("audit time is required")
	}
	copyDetail := map[string]string{}
	for key, value := range detail {
		copyDetail[key] = value
	}
	return model.AuditEntry{ID: id, CaseID: caseID, ActorID: actorID, Action: action, TargetID: targetID, Detail: copyDetail, OccurredAt: at}, nil
}

func Timeline(entries []model.AuditEntry, caseID string, limit int) []model.ActivitySummary {
	if limit <= 0 {
		limit = 20
	}
	result := make([]model.ActivitySummary, 0, limit)
	for i := len(entries) - 1; i >= 0 && len(result) < limit; i-- {
		entry := entries[i]
		if entry.CaseID != caseID {
			continue
		}
		result = append(result, model.ActivitySummary{Action: entry.Action, ActorID: entry.ActorID, TargetID: entry.TargetID, OccurredAt: entry.OccurredAt})
	}
	return result
}
