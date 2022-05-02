package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/MarkoGlamocak/fetchRewardsAssessment/models"
	"net/http"
)

func GetTransactionRecords(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, models.TransactionRecords)
}
