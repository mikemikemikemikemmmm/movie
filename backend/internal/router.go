package internal

import (
	"backend/internal/config"
	"backend/internal/otel"
	"backend/internal/promethus"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func registerMiddleware(r *gin.Engine) {

	corsConfig := cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"},
	}
	r.Use(cors.New(corsConfig))
	r.Use(gin.Recovery())
	r.Use(promethus.TimerMiddleware())
	r.Use(otel.TracingMiddleware())
}
func registerRoutes(r *gin.Engine) {

	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(promethus.Registry, promhttp.HandlerOpts{})))
	r.GET("/seats", handleGetAllSeats)
	r.GET("/refresh_seats", handleRefreshSeats)
	r.POST("/check_reserve", handleCheckReserve)
	r.POST("/reserve", handleReserve)
	r.POST("/health", handleCheckHealth)
	r.POST("/ready", handleCheckReady)
}
func InitRouter() *http.Server {
	RootRouter := gin.Default()
	registerMiddleware(RootRouter)
	registerRoutes(RootRouter)
	apiServer := &http.Server{
		Addr:    ":" + config.GetConfig().Port,
		Handler: RootRouter,
	}
	return apiServer
}
