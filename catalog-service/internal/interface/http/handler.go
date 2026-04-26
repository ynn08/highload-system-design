package http

import (
	"github.com/gin-gonic/gin"
	"github.com/user/highload-system-design/catalog-service/internal/usecase"
	"net/http"
	"strconv"
)

type RestaurantHandler struct {
	searchUseCase  *usecase.SearchRestaurantsUseCase
	getMenuUseCase *usecase.GetRestaurantMenuUseCase
}

func NewRestaurantHandler(searchUseCase *usecase.SearchRestaurantsUseCase, getMenuUseCase *usecase.GetRestaurantMenuUseCase) *RestaurantHandler {
	return &RestaurantHandler{searchUseCase: searchUseCase, getMenuUseCase: getMenuUseCase}
}

func (h *RestaurantHandler) Search(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	radiusStr := c.DefaultQuery("radius", "5000")
	cuisine := c.Query("cuisine")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lon"})
		return
	}

	radius, _ := strconv.Atoi(radiusStr)

	output, err := h.searchUseCase.Execute(c.Request.Context(), usecase.SearchInput{
		Lat:     lat,
		Lon:     lon,
		Radius:  radius,
		Cuisine: cuisine,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  output.Items,
		"total":  output.Total,
		"limit":  10,
		"offset": 0,
	})
}

func (h *RestaurantHandler) GetMenu(c *gin.Context) {
	id := c.Param("id")
	menu, err := h.getMenuUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "restaurant not found"})
		return
	}
	c.JSON(http.StatusOK, menu)
}
