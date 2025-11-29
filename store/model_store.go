package store

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"anomaly_detector/models"
)

type IModelStore interface {
	Store(ctx context.Context, model *models.APIModel) error
	StoreAll(ctx context.Context, models []*models.APIModel) (int, error)
	Get(ctx context.Context, path, method string) (*models.APIModel, error)
}

type ModelStore struct {
	mu     sync.RWMutex
	models map[string]*models.APIModel
}

func NewModelStore() IModelStore {
	return &ModelStore{
		models: make(map[string]*models.APIModel),
	}
}

func (s *ModelStore) Store(ctx context.Context, model *models.APIModel) error {
	if model == nil {
		return fmt.Errorf("model cannot be nil")
	}

	if model.Path == "" || model.Method == "" {
		return fmt.Errorf("path and method are required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := models.Key(model.Path, model.Method)

	s.models[key] = model

	slog.InfoContext(ctx, "Model stored", "path", model.Path, "method", model.Method)

	return nil
}

func (s *ModelStore) StoreAll(ctx context.Context, models []*models.APIModel) (int, error) {
	storedCount := 0

	for _, model := range models {
		if err := s.Store(ctx, model); err != nil {
			continue
		}

		storedCount++
	}

	return storedCount, nil
}

func (s *ModelStore) Get(ctx context.Context, path, method string) (*models.APIModel, error) {
	slog.InfoContext(ctx, "Getting model", "path", path, "method", method)

	s.mu.RLock()
	defer s.mu.RUnlock()

	key := models.Key(path, method)

	model, exists := s.models[key]
	if !exists {
		return nil, fmt.Errorf("model not found for %s %s", method, path)
	}

	slog.InfoContext(ctx, "Model retrieved", "path", path, "method", method)

	return model, nil
}
