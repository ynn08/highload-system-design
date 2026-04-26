package usecase

import (
	"context"
	"github.com/user/highload-system-design/catalog-service/internal/domain"
)

type GetRestaurantMenuUseCase struct {
	repo domain.RestaurantRepository
}

func NewGetRestaurantMenuUseCase(repo domain.RestaurantRepository) *GetRestaurantMenuUseCase {
	return &GetRestaurantMenuUseCase{repo: repo}
}

func (u *GetRestaurantMenuUseCase) Execute(ctx context.Context, restaurantID string) ([]domain.MenuItem, error) {
	restaurant, err := u.repo.GetByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	return restaurant.Menu, nil
}
