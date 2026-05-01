package usecase

import (
	"context"
	"github.com/user/highload-system-design/catalog-service/internal/domain"
)

type SaveRestaurantUseCase struct {
	repo domain.RestaurantRepository
}

func NewSaveRestaurantUseCase(repo domain.RestaurantRepository) *SaveRestaurantUseCase {
	return &SaveRestaurantUseCase{repo: repo}
}

func (u *SaveRestaurantUseCase) Execute(ctx context.Context, restaurant domain.Restaurant) error {
	return u.repo.Save(ctx, restaurant)
}
