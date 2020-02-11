package usecase

import (
	"context"
	"github.com/sschiz/apartament/internal/apartment"
	"github.com/sschiz/apartament/models"
)

type ApartmentUseCase struct{
	aptRepo apartment.Repository
}

func (uc ApartmentUseCase) Create(ctx context.Context, apartment *models.Apartment) error {
	return uc.aptRepo.Create(ctx, apartment)
}

func (uc ApartmentUseCase) Get(ctx context.Context, apartment *models.Apartment, opts ...apartment.Option) ([]*models.Apartment, error) {
	return uc.aptRepo.Get(ctx, apartment, opts...)
}

