package service

import "fmt"

type Resource interface{ Release() }
type ResourcePool interface {
	Acquire(documentID string) (Resource, error)
}

func (s *Service) ExportBatch(documentIDs []string, pool ResourcePool) ([][]byte, error) {
	result := make([][]byte, 0, len(documentIDs))
	for _, documentID := range documentIDs {
		resource, err := pool.Acquire(documentID)
		if err != nil {
			return nil, fmt.Errorf("acquire export resource for %s: %w", documentID, err)
		}
		version, err := s.store.CurrentVersion(documentID)
		if err != nil {
			resource.Release()
			return nil, err
		}
		result = append(result, append([]byte(nil), version.Content...))
		// Each record's resource is released as soon as its work is done so
		// the pool's quota is not held until the whole batch finishes.
		resource.Release()
	}
	return result, nil
}
