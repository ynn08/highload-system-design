package usecase

import (
	"context"
	"github.com/user/highload-system-design/delivery-service/internal/domain"
)

type SaveCourierUseCase struct {
	repo domain.CourierRepository
}

func NewSaveCourierUseCase(repo domain.CourierRepository) *SaveCourierUseCase {
	return &SaveCourierUseCase{repo: repo}
}

func (u *SaveCourierUseCase) Execute(ctx context.Context, courier *domain.Courier) error {
	return u.repo.SaveCourier(ctx, courier)
}
