package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	bolt "go.etcd.io/bbolt"
)

var (
	casesBucket     = []byte("cases")
	foldersBucket   = []byte("folders")
	documentsBucket = []byte("documents")
	versionsBucket  = []byte("versions")
	grantsBucket    = []byte("grants")
	sharesBucket    = []byte("shares")
	auditBucket     = []byte("audit")
)

type Store struct {
	db   *bolt.DB
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	store := &Store{db: db, path: path}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, name := range [][]byte{casesBucket, foldersBucket, documentsBucket, versionsBucket, grantsBucket, sharesBucket, auditBucket} {
			if _, createErr := tx.CreateBucketIfNotExists(name); createErr != nil {
				return createErr
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	return store, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func putJSON(bucket *bolt.Bucket, key string, value any) error {
	if key == "" {
		return errors.New("storage key is required")
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode %s: %w", key, err)
	}
	return bucket.Put([]byte(key), encoded)
}

func getJSON(bucket *bolt.Bucket, key string, destination any) error {
	encoded := bucket.Get([]byte(key))
	if encoded == nil {
		return fmt.Errorf("%s: %w", key, ErrNotFound)
	}
	if err := json.Unmarshal(encoded, destination); err != nil {
		return fmt.Errorf("decode %s: %w", key, err)
	}
	return nil
}

func listJSON[T any](bucket *bolt.Bucket) ([]T, error) {
	values := []T{}
	err := bucket.ForEach(func(_, encoded []byte) error {
		var value T
		if err := json.Unmarshal(encoded, &value); err != nil {
			return err
		}
		values = append(values, value)
		return nil
	})
	return values, err
}

var ErrNotFound = errors.New("not found")
