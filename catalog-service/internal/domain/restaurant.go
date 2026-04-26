package domain

import "context"

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type MenuItem struct {
	ID          string  `json:"item_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
}

type Restaurant struct {
	ID             string     `json:"restaurant_id"`
	Name           string     `json:"name"`
	Location       GeoPoint   `json:"location"`
	Cuisines       []string   `json:"cuisines"`
	Rating         float64    `json:"rating"`
	DistanceMeters int        `json:"distance_meters,omitempty"`
	Score          float64    `json:"-"` // Internal ranking score
	Menu           []MenuItem `json:"menu,omitempty"`
}

type RestaurantRepository interface {
	Search(ctx context.Context, lat, lon float64, radius int, cuisine string) ([]Restaurant, int, error)
	GetByID(ctx context.Context, id string) (*Restaurant, error)
	Save(ctx context.Context, restaurant Restaurant) error
}
