package internal

import (
	"backend/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var RootRouter *gin.Engine

func InitRouter() {
	RootRouter = gin.Default()
	corsConfig := cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"},
	}
	RootRouter.Use(cors.New(corsConfig))
	RootRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	RootRouter.GET("/seats", handleGetAllSeats)
	RootRouter.GET("/refresh_seats", handleRefreshSeats)
	RootRouter.POST("/check_reserve", handleCheckReserve)
	RootRouter.POST("/reserve", handleReserve)
	RootRouter.Run("0.0.0.0:" + config.GetConfig().Port)
}
