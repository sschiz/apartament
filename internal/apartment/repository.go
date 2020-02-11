package apartment

import (
	"context"
	"github.com/sschiz/apartament/models"
)

// Repository is apartment repository interface
type Repository interface {
	Create(ctx context.Context, apartment *models.Apartment) error
	Get(ctx context.Context, apartment *models.Apartment, opts ...Option) ([]*models.Apartment, error)
}
