package usecase

import (
	"context"
	"github.com/sschiz/apartament/internal/apartment"
	"github.com/sschiz/apartament/models"
)

// ApartmentUseCase is apartment.UseCase implementing
type ApartmentUseCase struct {
	aptRepo apartment.Repository
}

// NewApartmentUseCase create use case for apartment
func NewApartmentUseCase(aptRepo apartment.Repository) *ApartmentUseCase {
	return &ApartmentUseCase{aptRepo: aptRepo}
}

// Create creates apartment
func (uc ApartmentUseCase) Create(ctx context.Context, apartment *models.Apartment) error {
	return uc.aptRepo.Create(ctx, apartment)
}

// Get return apartments by apartment model
func (uc ApartmentUseCase) Get(ctx context.Context, apartment *models.Apartment, opts ...apartment.Option) ([]*models.Apartment, error) {
	return uc.aptRepo.Get(ctx, apartment, opts...)
}
