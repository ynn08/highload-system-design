package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	r := gin.Default()

	limiter := rate.NewLimiter(100, 200)
	r.Use(rateLimitMiddleware(limiter))

	catalogURL := os.Getenv("CATALOG_SERVICE_URL")
	orderURL := os.Getenv("ORDER_SERVICE_URL")
	deliveryURL := os.Getenv("DELIVERY_SERVICE_URL")

	r.Any("/api/v1/*any", func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// Order Status
		if strings.HasSuffix(path, "/status") && strings.Contains(path, "/orders/") {
			proxy(c, deliveryURL)
			return
		}

		// Carts
		if strings.HasPrefix(path, "/api/v1/carts") {
			proxy(c, orderURL)
			return
		}

		// Orders
		if strings.HasPrefix(path, "/api/v1/orders") {
			proxy(c, orderURL)
			return
		}

		// Restaurants
		if strings.HasPrefix(path, "/api/v1/restaurants") {
			proxy(c, catalogURL)
			return
		}

		// Couriers
		if strings.HasPrefix(path, "/api/v1/couriers") {
			proxy(c, deliveryURL)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "route not found in gateway"})
	})

	r.Any("/ws/*any", func(c *gin.Context) {
		proxy(c, deliveryURL)
	})

	r.GET("/health", func(c *gin.Context) { c.Status(200) })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on :%s", port)
	r.Run(":" + port)
}

func proxy(c *gin.Context, target string) {
	if target == "" {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	remote, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func rateLimitMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
