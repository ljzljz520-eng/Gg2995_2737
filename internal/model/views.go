package model

type DocumentSummary struct {
	ID, Name, Kind, Status                     string
	Version                                    int
	CanPreview, CanEdit, CanDownload, CanShare bool
}
type CaseSummary struct {
	ID, Number, Title, Client, Status string
	DocumentCount, FolderCount        int
}
type ActivitySummary struct {
	Action, ActorID, TargetID string
	OccurredAt                int64
}
type SearchFilter struct {
	CaseID, Query, Kind, Status string
	Limit, Offset               int
}
type Dashboard struct{ OpenCases, StayedCases, ClosedCases, Documents, RecentEdits, ActiveShares int }

func NormalizeFilter(f SearchFilter) SearchFilter {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Limit > 200 {
		f.Limit = 200
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return f
}
