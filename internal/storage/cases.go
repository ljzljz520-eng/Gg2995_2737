package storage

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"lawdrive/internal/model"
)

func (s *Store) PutCase(value model.CaseMatter) error {
	if err := value.Validate(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(casesBucket)
		if bucket.Get([]byte(value.ID)) != nil {
			return fmt.Errorf("case %s already exists", value.ID)
		}
		return putJSON(bucket, value.ID, value)
	})
}

func (s *Store) UpdateCase(value model.CaseMatter) error {
	if err := value.Validate(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(casesBucket)
		if bucket.Get([]byte(value.ID)) == nil {
			return fmt.Errorf("case %s: %w", value.ID, ErrNotFound)
		}
		return putJSON(bucket, value.ID, value)
	})
}

func (s *Store) Case(id string) (model.CaseMatter, error) {
	var value model.CaseMatter
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(casesBucket), id, &value) })
	return value, err
}

func (s *Store) Cases() ([]model.CaseMatter, error) {
	var values []model.CaseMatter
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		values, err = listJSON[model.CaseMatter](tx.Bucket(casesBucket))
		return err
	})
	return values, err
}

func (s *Store) PutFolder(value model.Folder) error {
	if value.ID == "" || value.CaseID == "" || value.Name == "" {
		return fmt.Errorf("folder identity is incomplete")
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(casesBucket).Get([]byte(value.CaseID)) == nil {
			return fmt.Errorf("case %s: %w", value.CaseID, ErrNotFound)
		}
		return putJSON(tx.Bucket(foldersBucket), value.ID, value)
	})
}

func (s *Store) Folders(caseID string) ([]model.Folder, error) {
	all := []model.Folder{}
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		all, err = listJSON[model.Folder](tx.Bucket(foldersBucket))
		return err
	})
	if err != nil {
		return nil, err
	}
	result := []model.Folder{}
	for _, folder := range all {
		if folder.CaseID == caseID {
			result = append(result, folder)
		}
	}
	return result, nil
}
