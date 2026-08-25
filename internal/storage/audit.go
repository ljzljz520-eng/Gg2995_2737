package storage

import (
	bolt "go.etcd.io/bbolt"
	"lawdrive/internal/model"
)

func (s *Store) AppendAudit(value model.AuditEntry) error {
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(auditBucket), value.ID, value) })
}
func (s *Store) AuditEntries(caseID string) ([]model.AuditEntry, error) {
	all := []model.AuditEntry{}
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		all, err = listJSON[model.AuditEntry](tx.Bucket(auditBucket))
		return err
	})
	if err != nil {
		return nil, err
	}
	result := []model.AuditEntry{}
	for _, entry := range all {
		if caseID == "" || entry.CaseID == caseID {
			result = append(result, entry)
		}
	}
	return result, nil
}
