package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// DecodeCursor decodes all elements from the cursor into the specified slice
func DecodeCursor[T any](ctx context.Context, cursor *mongo.Cursor, logger *zap.Logger) ([]T, error) {
	var items []T
	for cursor.Next(ctx) {
		var item T
		if decodeErr := cursor.Decode(&item); decodeErr != nil {
			logger.Error("Failed to decode item", zap.Error(decodeErr))
			return nil, decodeErr
		}
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		logger.Error("Cursor encountered an error", zap.Error(err))
		return nil, err
	}

	return items, nil
}
