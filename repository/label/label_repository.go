package repository

import (
	"context"

	"github.com/geekakili/portside/models"
)

// LabelRepository ..
type LabelRepository interface {
	AddLabel(ctx context.Context, label string, description string) (*models.Label, error)
	GetLabel(ctx context.Context, label string) (models.Label, error)
	GetLabels(ctx context.Context) []models.Label
	Updatelabel(ctx context.Context, label string, labelData models.Label) error
	Delete(ctx context.Context, label string) (bool, error)
}
