package usecase

import (
	"context"
	"errors"
	"log"
	"math"

	"github.com/user/highload-system-design/delivery-service/internal/domain"
)

type CourierRepository interface {
	FindAll(ctx context.Context) ([]*domain.Courier, error)
	Update(ctx context.Context, courier *domain.Courier) error
}

type DeliveryRepository interface {
	Save(ctx context.Context, delivery *domain.Delivery) error
	FindByID(ctx context.Context, orderID string) (*domain.Delivery, error)
}

type AssignCourierUseCase struct {
	courierRepo  CourierRepository
	deliveryRepo DeliveryRepository
}

func NewAssignCourierUseCase(cr CourierRepository, dr DeliveryRepository) *AssignCourierUseCase {
	return &AssignCourierUseCase{courierRepo: cr, deliveryRepo: dr}
}

func (u *AssignCourierUseCase) DeliveryRepo() DeliveryRepository {
	return u.deliveryRepo
}

func (u *AssignCourierUseCase) Execute(ctx context.Context, orderID, restaurantID string, dest domain.Location) error {
	log.Printf("Assigning courier for order %s...", orderID)

	couriers, err := u.courierRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	var bestCourier *domain.Courier
	minDist := math.MaxFloat64

	// Courier batching strategy
	for _, c := range couriers {
		if !c.IsAvailable || len(c.CurrentOrders) >= 3 {
			continue
		}

		dist := calculateDistance(c.Location, dest)
		
		if len(c.CurrentOrders) > 0 {
			dist = dist * 0.7 
		}

		if dist < minDist {
			minDist = dist
			bestCourier = c
		}
	}

	if bestCourier == nil {
		return errors.New("no available couriers")
	}

	delivery := &domain.Delivery{
		OrderID:      orderID,
		RestaurantID: restaurantID,
		Status:       domain.StatusAssigned,
		CourierID:    bestCourier.ID,
		Destination:  dest,
	}

	if err := u.deliveryRepo.Save(ctx, delivery); err != nil {
		return err
	}

	bestCourier.CurrentOrders = append(bestCourier.CurrentOrders, orderID)
	return u.courierRepo.Update(ctx, bestCourier)
}

func calculateDistance(l1, l2 domain.Location) float64 {
	return math.Sqrt(math.Pow(l1.Lat-l2.Lat, 2) + math.Pow(l1.Lon-l2.Lon, 2))
}
