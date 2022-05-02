package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/MarkoGlamocak/fetchRewardsAssessment/handlers/helpers"
	"github.com/MarkoGlamocak/fetchRewardsAssessment/models"
	"net/http"
	"sort"
)

func AddTransactionRecord(context *gin.Context) {
	var newTransaction models.Transaction

	if err := context.BindJSON(&newTransaction); err != nil {
		return
	}

	pointBalances := helpers.CalculatePointBalances()

	if _, exists := pointBalances[newTransaction.Payer]; exists {
		if newTransaction.Points < 0 && pointBalances[newTransaction.Payer] + newTransaction.Points < 0 {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Payer Points Can't Be Negative"})
			return
		}
	} else if newTransaction.Points < 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Payer Points Can't Be Negative"})
		return
	}

	models.TransactionRecords = append(models.TransactionRecords, newTransaction)

	models.TotalPoints+=newTransaction.Points

	sort.Slice(models.TransactionRecords, func(i, j int) bool {
		return models.TransactionRecords[i].TimeStamp.Before(models.TransactionRecords[j].TimeStamp)
	})

	context.IndentedJSON(http.StatusCreated, newTransaction)
}
