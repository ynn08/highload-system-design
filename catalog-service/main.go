package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/user/highload-system-design/catalog-service/internal/domain"
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

	repo := elasticsearch.NewRestaurantRepository(client, domain.ElasticIndexRestaurants)
	
	searchUC := usecase.NewSearchRestaurantsUseCase(repo)
	getMenuUC := usecase.NewGetRestaurantMenuUseCase(repo)
	saveUC := usecase.NewSaveRestaurantUseCase(repo)
	handler := http.NewRestaurantHandler(searchUC, getMenuUC, saveUC)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/restaurants", handler.Search)
		v1.GET("/restaurants/:id/menu", handler.GetMenu)
		v1.POST("/restaurants", handler.Save) // Internal/Seeding endpoint
	}
	
	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	log.Println("Catalog service starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
