package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/user/highload-system-design/catalog-service/internal/domain"
)

type RestaurantRepository struct {
	client *elastic.Client
	index  string
}

func NewRestaurantRepository(client *elastic.Client, index string) *RestaurantRepository {
	return &RestaurantRepository{client: client, index: index}
}

func (r *RestaurantRepository) Search(ctx context.Context, lat, lon float64, radius int, cuisine string) ([]domain.Restaurant, int, error) {
	query := elastic.NewBoolQuery()

	geoQuery := elastic.NewGeoDistanceQuery("location").
		Lat(lat).
		Lon(lon).
		Distance(fmt.Sprintf("%dm", radius))
	query = query.Filter(geoQuery)

	if cuisine != "" {
		query = query.Must(elastic.NewMatchQuery("cuisines", cuisine))
	}

	searchResult, err := r.client.Search().
		Index(r.index).
		Query(query).
		Size(50).
		Do(ctx)

	if err != nil {
		return nil, 0, err
	}

	restaurants := make([]domain.Restaurant, 0)
	for _, hit := range searchResult.Hits.Hits {
		var res domain.Restaurant
		if err := json.Unmarshal(hit.Source, &res); err == nil {
			restaurants = append(restaurants, res)
		}
	}

	return restaurants, int(searchResult.TotalHits()), nil
}

func (r *RestaurantRepository) GetByID(ctx context.Context, id string) (*domain.Restaurant, error) {
	res, err := r.client.Get().
		Index(r.index).
		Id(id).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var restaurant domain.Restaurant
	if err := json.Unmarshal(res.Source, &restaurant); err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *RestaurantRepository) Save(ctx context.Context, restaurant domain.Restaurant) error {
	_, err := r.client.Index().
		Index(r.index).
		Id(restaurant.ID).
		BodyJson(restaurant).
		Do(ctx)
	return err
}
