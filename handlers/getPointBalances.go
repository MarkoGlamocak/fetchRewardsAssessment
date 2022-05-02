package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/MarkoGlamocak/fetchRewardsAssessment/handlers/helpers"
	"net/http"
)

func GetPointBalances(context *gin.Context) {
	pointBalances := helpers.CalculatePointBalances()

	context.IndentedJSON(http.StatusOK, pointBalances)
}
