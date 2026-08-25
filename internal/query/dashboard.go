package query

import "lawdrive/internal/model"

func (q *Queries) Dashboard(userID string, now int64) (model.Dashboard, error) {
	cases, err := q.store.Cases()
	if err != nil {
		return model.Dashboard{}, err
	}
	shares, err := q.store.Shares()
	if err != nil {
		return model.Dashboard{}, err
	}
	result := model.Dashboard{}
	for _, caseMatter := range cases {
		if _, ok := model.RoleFor(caseMatter, userID); !ok {
			continue
		}
		switch caseMatter.Status {
		case model.CaseOpen:
			result.OpenCases++
		case model.CaseStayed:
			result.StayedCases++
		case model.CaseClosed:
			result.ClosedCases++
		}
		documents, loadErr := q.store.Documents(caseMatter.ID)
		if loadErr != nil {
			return model.Dashboard{}, loadErr
		}
		result.Documents += len(documents)
	}
	for _, share := range shares {
		if !share.Revoked && share.ExpiresAt > now {
			result.ActiveShares++
		}
	}
	return result, nil
}
