package storage

import (
	bolt "go.etcd.io/bbolt"
	"lawdrive/internal/model"
)

func (s *Store) PutGrant(value model.PermissionGrant) error {
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(grantsBucket), value.ID, value) })
}
func (s *Store) Grants(documentID, subjectID string) ([]model.PermissionGrant, error) {
	all := []model.PermissionGrant{}
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		all, err = listJSON[model.PermissionGrant](tx.Bucket(grantsBucket))
		return err
	})
	result := []model.PermissionGrant{}
	for _, grant := range all {
		if grant.DocumentID == documentID && grant.SubjectID == subjectID {
			result = append(result, grant)
		}
	}
	return result, err
}
func (s *Store) PutShare(value model.ShareLink) error {
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(sharesBucket), value.ID, value) })
}
func (s *Store) ShareByToken(token string) (model.ShareLink, error) {
	all := []model.ShareLink{}
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		all, err = listJSON[model.ShareLink](tx.Bucket(sharesBucket))
		return err
	})
	if err != nil {
		return model.ShareLink{}, err
	}
	for _, share := range all {
		if share.Token == token {
			return share, nil
		}
	}
	return model.ShareLink{}, ErrNotFound
}
func (s *Store) Shares() ([]model.ShareLink, error) {
	var result []model.ShareLink
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		result, err = listJSON[model.ShareLink](tx.Bucket(sharesBucket))
		return err
	})
	return result, err
}
