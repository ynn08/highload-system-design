package usecase

import (
	"context"
	"github.com/user/highload-system-design/catalog-service/internal/domain"
	"testing"
)

type mockRepo struct {
	restaurants []domain.Restaurant
	total       int
	err         error
}

func (m *mockRepo) Search(ctx context.Context, lat, lon float64, radius int, cuisine string) ([]domain.Restaurant, int, error) {
	return m.restaurants, m.total, m.err
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*domain.Restaurant, error) {
	for _, r := range m.restaurants {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) Save(ctx context.Context, restaurant domain.Restaurant) error {
	return nil
}

func TestSearchRestaurantsUseCase_Execute(t *testing.T) {
	repo := &mockRepo{
		restaurants: []domain.Restaurant{
			{
				ID:       "1",
				Name:     "Far Restaurant",
				Location: domain.GeoPoint{Lat: 56.0, Lon: 38.0},
				Rating:   5.0,
			},
			{
				ID:       "2",
				Name:     "Near Restaurant",
				Location: domain.GeoPoint{Lat: 55.75, Lon: 37.62},
				Rating:   4.0,
			},
		},
		total: 2,
	}

	uc := NewSearchRestaurantsUseCase(repo)

	// User at 55.75, 37.62 (same as Near Restaurant)
	output, err := uc.Execute(context.Background(), SearchInput{
		Lat:    55.75,
		Lon:    37.62,
		Radius: 500000, // Large radius to include both
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(output.Items))
	}

	// Near Restaurant should be first due to distance penalty on the Far one
	if output.Items[0].ID != "2" {
		t.Errorf("expected first restaurant to be '2', got %s", output.Items[0].ID)
	}

	if output.Items[0].DistanceMeters != 0 {
		t.Errorf("expected distance 0 for nearest, got %d", output.Items[0].DistanceMeters)
	}
}
