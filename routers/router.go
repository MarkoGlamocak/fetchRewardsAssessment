package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/MarkoGlamocak/fetchRewardsAssessment/handlers"
)

func Init() *gin.Engine {
	router := gin.Default()
	router.GET("/transaction", handlers.GetTransactionRecords) // Route designed for testing purposes
	router.POST("/transaction", handlers.AddTransactionRecord)
	router.PUT("/points", handlers.SpendPoints)
	router.GET("/points", handlers.GetPointBalances)

	return router
}
