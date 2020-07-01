package repository

import (
	"context"

	"github.com/geekakili/portside/models"
)

// ImageRepository runs queries for images against the database
type ImageRepository interface {
	AddLabel(ctx context.Context, label string, imageTags ...string) error
	GetImageLabels(ctx context.Context, imageName string) (labels []string, err error)
	GetByLabel(ctx context.Context, label string) ([]*models.Image, error)
}
