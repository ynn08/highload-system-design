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

	// Rate Limiter: 100 requests per second with burst of 200
	limiter := rate.NewLimiter(100, 200)

	r.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	})

	// Auth Middleware (Mock)
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/orders") || strings.HasPrefix(c.Request.URL.Path, "/api/v1/carts") {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				// For demo purposes, we allow it but log a warning. 
				// In production, we would c.AbortWithStatus(401)
				log.Println("Warning: Missing Authorization header for protected route")
			}
		}
		c.Next()
	})

	// Reverse Proxy mapping
	targets := map[string]string{
		"/api/v1/restaurants": os.Getenv("CATALOG_SERVICE_URL"),
		"/api/v1/orders":      os.Getenv("ORDER_SERVICE_URL"),
		"/api/v1/carts":       os.Getenv("ORDER_SERVICE_URL"),
		"/ws":                 os.Getenv("DELIVERY_SERVICE_URL"),
	}

	for path, target := range targets {
		if target == "" {
			continue
		}
		proxyPath := path
		targetURL, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		r.Any(proxyPath+"/*any", func(c *gin.Context) {
			proxy.ServeHTTP(c.Writer, c.Request)
		})
		
		// Also handle exact path
		r.Any(proxyPath, func(c *gin.Context) {
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on :%s", port)
	r.Run(":" + port)
}
