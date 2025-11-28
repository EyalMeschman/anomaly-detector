package store

import (
	"context"
	"testing"

	"anomaly_detector/models"

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
		err := store.Store(ctx, tModel)
		assert.NoError(t, err)

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

func TestInMemoryModelStore_StoreAll(t *testing.T) {
	t.Run("success storing multiple models", func(t *testing.T) {
		ctx := context.Background()
		tStore := NewModelStore()

		tModels := []*models.APIModel{
			{Path: "/users", Method: "GET"},
			{Path: "/users", Method: "POST"},
			{Path: "/products", Method: "GET"},
		}

		tCount, tErr := tStore.StoreAll(ctx, tModels)

		assert.NoError(t, tErr)
		assert.Equal(t, 3, tCount)

		// Verify all models were stored by retrieving them individually
		_, err := tStore.Get(ctx, "/users", "GET")
		assert.NoError(t, err)

		_, err = tStore.Get(ctx, "/users", "POST")
		assert.NoError(t, err)

		_, err = tStore.Get(ctx, "/products", "GET")
		assert.NoError(t, err)
	})

	t.Run("partial success with invalid models", func(t *testing.T) {
		ctx := context.Background()
		tStore := NewModelStore()

		tModels := []*models.APIModel{
			{Path: "/users", Method: "GET"},
			{Path: "", Method: "POST"}, // Invalid - empty path
			{Path: "/products", Method: "GET"},
		}

		tCount, tErr := tStore.StoreAll(ctx, tModels)

		assert.NoError(t, tErr)
		assert.Equal(t, 2, tCount) // Only 2 valid models stored

		// Verify only valid models were stored
		_, err := tStore.Get(ctx, "/users", "GET")
		assert.NoError(t, err)

		_, err = tStore.Get(ctx, "/products", "GET")
		assert.NoError(t, err)
		// Invalid model should not be stored
		_, err = tStore.Get(ctx, "", "POST")
		assert.Error(t, err)
	})

	t.Run("empty slice returns zero", func(t *testing.T) {
		ctx := context.Background()
		tStore := NewModelStore()

		tModels := []*models.APIModel{}

		tCount, tErr := tStore.StoreAll(ctx, tModels)

		assert.NoError(t, tErr)
		assert.Equal(t, 0, tCount)
	})
}
