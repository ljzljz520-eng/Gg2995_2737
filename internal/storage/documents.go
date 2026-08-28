package storage

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"lawdrive/internal/model"
)

func (s *Store) CreateDocument(document model.Document, version model.DocumentVersion) error {
	if err := document.Validate(); err != nil {
		return err
	}
	if err := version.Validate(); err != nil {
		return err
	}
	if version.DocumentID != document.ID || version.Sequence != document.CurrentVersion {
		return fmt.Errorf("initial version does not match document")
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(casesBucket).Get([]byte(document.CaseID)) == nil {
			return fmt.Errorf("case %s: %w", document.CaseID, ErrNotFound)
		}
		if tx.Bucket(documentsBucket).Get([]byte(document.ID)) != nil {
			return fmt.Errorf("document %s already exists", document.ID)
		}
		if err := putJSON(tx.Bucket(documentsBucket), document.ID, document); err != nil {
			return err
		}
		return putJSON(tx.Bucket(versionsBucket), version.ID, version)
	})
}

func (s *Store) Document(id string) (model.Document, error) {
	var value model.Document
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(documentsBucket), id, &value) })
	return value, err
}

func (s *Store) Documents(caseID string) ([]model.Document, error) {
	all := []model.Document{}
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		all, err = listJSON[model.Document](tx.Bucket(documentsBucket))
		return err
	})
	if err != nil {
		return nil, err
	}
	result := []model.Document{}
	for _, document := range all {
		if document.CaseID == caseID {
			result = append(result, document)
		}
	}
	return result, nil
}

func (s *Store) AddVersion(version model.DocumentVersion) (model.Document, error) {
	if err := version.Validate(); err != nil {
		return model.Document{}, err
	}
	var updated model.Document
	err := s.db.Update(func(tx *bolt.Tx) error {
		if err := getJSON(tx.Bucket(documentsBucket), version.DocumentID, &updated); err != nil {
			return err
		}
		if version.Sequence != updated.CurrentVersion+1 {
			return fmt.Errorf("expected version %d", updated.CurrentVersion+1)
		}
		updated.CurrentVersion = version.Sequence
		if err := putJSON(tx.Bucket(versionsBucket), version.ID, version); err != nil {
			return err
		}
		return putJSON(tx.Bucket(documentsBucket), updated.ID, updated)
	})
	return updated, err
}

func (s *Store) Versions(documentID string) ([]model.DocumentVersion, error) {
	all := []model.DocumentVersion{}
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		all, err = listJSON[model.DocumentVersion](tx.Bucket(versionsBucket))
		return err
	})
	if err != nil {
		return nil, err
	}
	result := []model.DocumentVersion{}
	for _, version := range all {
		if version.DocumentID == documentID {
			result = append(result, version)
		}
	}
	return result, nil
}

func (s *Store) CurrentVersion(documentID string) (model.DocumentVersion, error) {
	document, err := s.Document(documentID)
	if err != nil {
		return model.DocumentVersion{}, err
	}
	versions, err := s.Versions(documentID)
	if err != nil {
		return model.DocumentVersion{}, err
	}
	for _, version := range versions {
		if version.Sequence == document.CurrentVersion {
			return version, nil
		}
	}
	return model.DocumentVersion{}, fmt.Errorf("current version for %s: %w", documentID, ErrNotFound)
}
