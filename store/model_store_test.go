package store

import (
	"context"
	"testing"

	"anomaly_detector/models"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("success getting existing model", func(t *testing.T) {
		ctx := context.Background()
		store := NewModelStore()
		tModel := &models.APIModel{
			Path:   "/users",
			Method: "POST",
		}

		// Store first using StoreAll
		err := store.StoreAll(ctx, []*models.APIModel{tModel})
		assert.NoError(t, err)

		// Get it back
		retrieved, err := store.Get(ctx, "/users", "POST")
		assert.NoError(t, err)
		assert.Equal(t, tModel, retrieved)
	})

	t.Run("error when model not found", func(t *testing.T) {
		ctx := context.Background()
		store := NewModelStore()

		retrieved, err := store.Get(ctx, "/nonexistent", "GET")
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})
}

func TestStoreAll(t *testing.T) {
	tModels := []*models.APIModel{
		{Path: "/users", Method: "GET"},
		{Path: "/users", Method: "POST"},
		{Path: "/products", Method: "GET"},
	}

	t.Run("success storing multiple models", func(t *testing.T) {
		ctx := context.Background()
		tStore := NewModelStore()

		tErr := tStore.StoreAll(ctx, tModels)

		assert.NoError(t, tErr)
	})

	t.Run("empty slice succeeds", func(t *testing.T) {
		ctx := context.Background()
		tStore := NewModelStore()

		tModels := []*models.APIModel{}

		err := tStore.StoreAll(ctx, tModels)
		assert.NoError(t, err)
	})

	t.Run("fail on nil model", func(t *testing.T) {
		ctx := context.Background()
		tStore := NewModelStore()

		invalidModels := tModels
		invalidModels = append(invalidModels, nil) // Nil model

		err := tStore.StoreAll(ctx, invalidModels)
		assert.Error(t, err)
	})

	t.Run("fail on empty path", func(t *testing.T) {
		ctx := context.Background()
		tStore := NewModelStore()

		invalidModels := tModels
		invalidModels = append(invalidModels, &models.APIModel{Path: "", Method: "POST"})

		err := tStore.StoreAll(ctx, invalidModels)
		assert.Error(t, err)
	})

	t.Run("fail on empty method", func(t *testing.T) {
		ctx := context.Background()
		tStore := NewModelStore()

		invalidModels := tModels
		invalidModels = append(invalidModels, &models.APIModel{Path: "path", Method: ""})

		err := tStore.StoreAll(ctx, invalidModels)
		assert.Error(t, err)
	})

	t.Run("fail on duplicate model", func(t *testing.T) {
		ctx := context.Background()
		tStore := NewModelStore()

		// Store first batch
		firstBatch := []*models.APIModel{
			{Path: "/users", Method: "GET"},
		}
		err := tStore.StoreAll(ctx, firstBatch)
		assert.NoError(t, err)

		// Try to store duplicate
		secondBatch := []*models.APIModel{
			{Path: "/products", Method: "GET"},
			{Path: "/users", Method: "GET"}, // Duplicate
		}
		err = tStore.StoreAll(ctx, secondBatch)
		assert.Error(t, err)
	})
}
