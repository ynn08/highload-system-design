package usecase

import (
	"context"
	"github.com/user/highload-system-design/order-service/internal/domain"
)

type ManageCartUseCase struct {
	repo domain.CartRepository
}

func NewManageCartUseCase(repo domain.CartRepository) *ManageCartUseCase {
	return &ManageCartUseCase{repo: repo}
}

func (u *ManageCartUseCase) AddToCart(ctx context.Context, customerID, restaurantID string, item domain.CartItem) (*domain.Cart, error) {
	cart, err := u.repo.FindByCustomerID(ctx, customerID)
	if err != nil || cart == nil {
		cart = &domain.Cart{CustomerID: customerID, RestaurantID: restaurantID}
	}

	if cart.RestaurantID != restaurantID {
		cart = &domain.Cart{CustomerID: customerID, RestaurantID: restaurantID}
	}

	cart.Items = append(cart.Items, item)
	err = u.repo.Save(ctx, cart)
	return cart, err
}

func (u *ManageCartUseCase) GetCart(ctx context.Context, customerID string) (*domain.Cart, error) {
	return u.repo.FindByCustomerID(ctx, customerID)
}
