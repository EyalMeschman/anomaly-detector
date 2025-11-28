package store

import (
	"context"
	"testing"

	"anomaly_detector/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryModelStore_Store(t *testing.T) {
	t.Run("success storing valid model", func(t *testing.T) {
		ctx := context.Background()
		store := NewModelStore()
		tModel := &models.APIModel{
			Path:   "/test",
			Method: "GET",
		}

		err := store.Store(ctx, tModel)
		assert.NoError(t, err)

		retrieved, err := store.Get(ctx, "/test", "GET")
		assert.NoError(t, err)
		assert.Equal(t, tModel, retrieved)
	})

	t.Run("error when path is empty", func(t *testing.T) {
		ctx := context.Background()
		store := NewModelStore()
		tModel := &models.APIModel{
			Path:   "",
			Method: "GET",
		}

		err := store.Store(ctx, tModel)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")
	})

	t.Run("error when method is empty", func(t *testing.T) {
		ctx := context.Background()
		store := NewModelStore()
		tModel := &models.APIModel{
			Path:   "/test",
			Method: "",
		}

		err := store.Store(ctx, tModel)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")
	})
}

func TestInMemoryModelStore_Get(t *testing.T) {
	t.Run("success getting existing model", func(t *testing.T) {
		ctx := context.Background()
		store := NewModelStore()
		tModel := &models.APIModel{
			Path:   "/users",
			Method: "POST",
		}

		// Store first
		store.Store(ctx, tModel)

		// Get it back
		retrieved, err := store.Get(ctx, "/users", "POST")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, tModel.Path, retrieved.Path)
		assert.Equal(t, tModel.Method, retrieved.Method)
	})

	t.Run("error when model not found", func(t *testing.T) {
		ctx := context.Background()
		store := NewModelStore()

		retrieved, err := store.Get(ctx, "/nonexistent", "GET")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestInMemoryModelStore_GetAll(t *testing.T) {
	t.Run("success getting all models", func(t *testing.T) {
		ctx := context.Background()
		store := NewModelStore()

		// Store 3 models
		tModel1 := &models.APIModel{Path: "/users", Method: "GET"}
		tModel2 := &models.APIModel{Path: "/users", Method: "POST"}
		tModel3 := &models.APIModel{Path: "/products", Method: "GET"}

		store.Store(ctx, tModel1)
		store.Store(ctx, tModel2)
		store.Store(ctx, tModel3)

		// Get all
		allModels := store.GetAll(ctx)
		assert.Len(t, allModels, 3)
	})

	t.Run("returns empty slice when no models", func(t *testing.T) {
		ctx := context.Background()
		store := NewModelStore()

		allModels := store.GetAll(ctx)
		assert.NotNil(t, allModels)
		assert.Len(t, allModels, 0)
	})
}
