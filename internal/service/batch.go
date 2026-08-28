package service

import "fmt"

type Resource interface{ Release() }
type ResourcePool interface {
	Acquire(documentID string) (Resource, error)
}

func (s *Service) ExportBatch(documentIDs []string, pool ResourcePool) ([][]byte, error) {
	resources := make([]Resource, 0, len(documentIDs))
	defer func() {
		for _, resource := range resources {
			resource.Release()
		}
	}()
	result := make([][]byte, 0, len(documentIDs))
	for _, documentID := range documentIDs {
		resource, err := pool.Acquire(documentID)
		if err != nil {
			return nil, fmt.Errorf("acquire export resource for %s: %w", documentID, err)
		}
		resources = append(resources, resource)
		version, err := s.store.CurrentVersion(documentID)
		if err != nil {
			return nil, err
		}
		result = append(result, append([]byte(nil), version.Content...))
	}
	return result, nil
}
