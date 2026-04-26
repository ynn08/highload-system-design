package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/user/highload-system-design/catalog-service/internal/infrastructure/elasticsearch"
	"github.com/user/highload-system-design/catalog-service/internal/interface/http"
	"github.com/user/highload-system-design/catalog-service/internal/usecase"
	"log"
	"os"
)

func main() {
	elasticURL := os.Getenv("ELASTICSEARCH_URL")
	if elasticURL == "" {
		elasticURL = "http://localhost:9200"
	}

	client, err := elastic.NewClient(
		elastic.SetURL(elasticURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		log.Fatalf("Error creating elastic client: %s", err)
	}

	repo := elasticsearch.NewRestaurantRepository(client, "restaurants")
	searchUseCase := usecase.NewSearchRestaurantsUseCase(repo)
	getMenuUseCase := usecase.NewGetRestaurantMenuUseCase(repo)
	handler := http.NewRestaurantHandler(searchUseCase, getMenuUseCase)

	r := gin.Default()
	r.GET("/api/v1/restaurants", handler.Search)
	r.GET("/api/v1/restaurants/:id/menu", handler.GetMenu)
	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	log.Println("Catalog service starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
