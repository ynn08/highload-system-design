package usecase

import (
	"context"
	"github.com/user/highload-system-design/catalog-service/internal/domain"
	"math"
	"sort"
)

type SearchRestaurantsUseCase struct {
	repo domain.RestaurantRepository
}

func NewSearchRestaurantsUseCase(repo domain.RestaurantRepository) *SearchRestaurantsUseCase {
	return &SearchRestaurantsUseCase{repo: repo}
}

type SearchInput struct {
	Lat     float64
	Lon     float64
	Radius  int
	Cuisine string
}

type SearchOutput struct {
	Items []domain.Restaurant
	Total int
}

func (u *SearchRestaurantsUseCase) Execute(ctx context.Context, input SearchInput) (SearchOutput, error) {
	restaurants, total, err := u.repo.Search(ctx, input.Lat, input.Lon, input.Radius, input.Cuisine)
	if err != nil {
		return SearchOutput{}, err
	}

	// Non-trivial logic: Ranking and Exact Distance Calculation
	for i := range restaurants {
		dist := haversine(input.Lat, input.Lon, restaurants[i].Location.Lat, restaurants[i].Location.Lon)
		restaurants[i].DistanceMeters = int(dist)
		
		// Ranking formula: Rating * 100 - (Distance / 100)
		// Higher rating increases score, further distance decreases it.
		restaurants[i].Score = (restaurants[i].Rating * 100) - (dist / 50)
	}

	sort.Slice(restaurants, func(i, j int) bool {
		return restaurants[i].Score > restaurants[j].Score
	})

	return SearchOutput{
		Items: restaurants,
		Total: total,
	}, nil
}

// haversine calculates distance between two points in meters
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
