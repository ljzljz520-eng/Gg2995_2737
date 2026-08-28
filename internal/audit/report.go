package audit

import (
	"lawdrive/internal/model"
	"sort"
)

type ActorReport struct {
	ActorID         string
	Actions         int
	FirstAt, LastAt int64
	ByAction        map[string]int
}

func Summarize(entries []model.AuditEntry) []ActorReport {
	byActor := map[string]*ActorReport{}
	for _, entry := range entries {
		report := byActor[entry.ActorID]
		if report == nil {
			report = &ActorReport{ActorID: entry.ActorID, FirstAt: entry.OccurredAt, ByAction: map[string]int{}}
			byActor[entry.ActorID] = report
		}
		report.Actions++
		report.ByAction[entry.Action]++
		if entry.OccurredAt < report.FirstAt {
			report.FirstAt = entry.OccurredAt
		}
		if entry.OccurredAt > report.LastAt {
			report.LastAt = entry.OccurredAt
		}
	}
	result := make([]ActorReport, 0, len(byActor))
	for _, report := range byActor {
		result = append(result, *report)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ActorID < result[j].ActorID })
	return result
}
