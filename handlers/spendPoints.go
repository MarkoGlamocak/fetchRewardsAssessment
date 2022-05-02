package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/MarkoGlamocak/fetchRewardsAssessment/handlers/helpers"
	"github.com/MarkoGlamocak/fetchRewardsAssessment/models"
	"net/http"
	"time"
)

func SpendPoints(context *gin.Context) {
	var pointsToSpend models.Points
	pointsSpentPerPayerMap := make(map[string]int)
	var pointsSpentPerPayer []models.PointsPerPayer

	if err := context.BindJSON(&pointsToSpend); err != nil {
		return
	}

	if pointsToSpend.Points > models.TotalPoints {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Can't spend more than you have"})
		return
	} else if pointsToSpend.Points < 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Can't spend negative points"})
		return
	}
	models.TotalPoints-=pointsToSpend.Points

	negativePointBalances := helpers.CalculateNegativePointBalances()

	for _, t := range models.TransactionRecords {
		if pointsToSpend.Points == 0 {
			break
		}

		if t.Points < 0 {
			continue
		}
		
		if _, exists := pointsSpentPerPayerMap[t.Payer]; exists {
			if negativePointBalances[t.Payer] == 0 {
				if t.Points > pointsToSpend.Points {
					pointsSpentPerPayerMap[t.Payer]-=pointsToSpend.Points
					pointsToSpend.Points = 0
				} else {
					pointsSpentPerPayerMap[t.Payer]-=t.Points
					pointsToSpend.Points-=t.Points
				}
			} else {
				if t.Points > -1*negativePointBalances[t.Payer] {
					pointsLeft := (t.Points + negativePointBalances[t.Payer])
					negativePointBalances[t.Payer] = 0
					if pointsLeft > pointsToSpend.Points {
						pointsSpentPerPayerMap[t.Payer]-=pointsToSpend.Points
						pointsToSpend.Points = 0
					} else {
						pointsSpentPerPayerMap[t.Payer]-=pointsLeft
						pointsToSpend.Points-=pointsLeft
					}
				} else {
					negativePointBalances[t.Payer]+=t.Points
				}
			}
		} else {
			if negativePointBalances[t.Payer] == 0 {
				if t.Points > pointsToSpend.Points {
					pointsSpentPerPayerMap[t.Payer] = pointsToSpend.Points * -1
					pointsToSpend.Points = 0
				} else {
					pointsSpentPerPayerMap[t.Payer] = t.Points * -1
					pointsToSpend.Points-=t.Points
				}
			} else {
				if t.Points > -1*negativePointBalances[t.Payer] {
					pointsLeft := (t.Points + negativePointBalances[t.Payer])
					negativePointBalances[t.Payer] = 0
					if pointsLeft > pointsToSpend.Points {
						pointsSpentPerPayerMap[t.Payer] = pointsToSpend.Points * -1
						pointsToSpend.Points = 0
					} else {
						pointsSpentPerPayerMap[t.Payer] = pointsLeft * -1
						pointsToSpend.Points-=pointsLeft
					}
				} else {
					negativePointBalances[t.Payer]+=t.Points
				}
			}
		}
	}

	for k, v := range pointsSpentPerPayerMap {
		pointsSpentPerPayer = append(pointsSpentPerPayer, models.PointsPerPayer{Payer: k, Points: v})
		models.TransactionRecords = append(models.TransactionRecords, models.Transaction{Payer: k, Points: v, TimeStamp: time.Now()})
	}

	context.IndentedJSON(http.StatusOK, pointsSpentPerPayer)
}
