package query

import (
	"lawdrive/internal/authz"
	"lawdrive/internal/model"
	"lawdrive/internal/storage"
	"strings"
)

type Queries struct{ store *storage.Store }

func New(store *storage.Store) *Queries { return &Queries{store: store} }

func (q *Queries) SearchDocuments(userID string, filter model.SearchFilter) ([]model.DocumentSummary, error) {
	filter = model.NormalizeFilter(filter)
	caseMatter, err := q.store.Case(filter.CaseID)
	if err != nil {
		return nil, err
	}
	role, member := model.RoleFor(caseMatter, userID)
	if !member {
		return []model.DocumentSummary{}, nil
	}
	documents, err := q.store.Documents(filter.CaseID)
	if err != nil {
		return nil, err
	}
	result := []model.DocumentSummary{}
	for _, document := range documents {
		if filter.Query != "" && !strings.Contains(strings.ToLower(document.Name), strings.ToLower(filter.Query)) {
			continue
		}
		if filter.Kind != "" && string(document.Kind) != filter.Kind {
			continue
		}
		if filter.Status != "" && document.Status != filter.Status {
			continue
		}
		if filter.Offset > 0 {
			filter.Offset--
			continue
		}
		result = append(result, model.DocumentSummary{ID: document.ID, Name: document.Name, Kind: string(document.Kind), Status: document.Status, Version: document.CurrentVersion, CanPreview: authz.Decide(role, authz.View, caseMatter.Status).Allowed, CanEdit: authz.Decide(role, authz.Edit, caseMatter.Status).Allowed, CanDownload: authz.Decide(role, authz.Download, caseMatter.Status).Allowed, CanShare: authz.Decide(role, authz.Share, caseMatter.Status).Allowed})
		if len(result) >= filter.Limit {
			break
		}
	}
	return result, nil
}
