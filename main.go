package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
)

type Transaction struct {
	Payer string `json:"payer"`
	Points int `json:"points"`
	TimeStamp time.Time `json:"timestamp"`
}

type Points struct {
	Points int `json:"points"`
}

type PointsPerPayer struct {
	Payer string `json:"payer"`
	Points int `json:"points"`
}

var totalPoints int = 0

var transactionRecords = []Transaction {}

// ROUTE HANDLERS

func getTransactionRecords(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, transactionRecords)
}

func addTransactionRecord(context *gin.Context) {
	var newTransaction Transaction

	if err := context.BindJSON(&newTransaction); err != nil {
		return
	}

	pointBalances := calculatePointBalances()

	if _, exists := pointBalances[newTransaction.Payer]; exists {
		if newTransaction.Points < 0 && pointBalances[newTransaction.Payer] + newTransaction.Points < 0 {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Payer Points Can't Be Negative"})
			return
		}
	} else if newTransaction.Points < 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Payer Points Can't Be Negative"})
		return
	}

	transactionRecords = append(transactionRecords, newTransaction)

	totalPoints+=newTransaction.Points

	sort.Slice(transactionRecords, func(i, j int) bool {
		return transactionRecords[i].TimeStamp.Before(transactionRecords[j].TimeStamp)
	})

	context.IndentedJSON(http.StatusCreated, newTransaction)
}

func spendPoints(context *gin.Context) {
	var pointsToSpend Points
	pointsSpentPerPayerMap := make(map[string]int)
	var pointsSpentPerPayer []PointsPerPayer

	if err := context.BindJSON(&pointsToSpend); err != nil {
		return
	}

	if pointsToSpend.Points > totalPoints {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Can't spend more than you have"})
		return
	} else if pointsToSpend.Points < 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Can't spend negative points"})
		return
	}
	totalPoints-=pointsToSpend.Points

	negativePointBalances := calculateNegativePointBalances()

	for _, t := range transactionRecords {
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
		pointsSpentPerPayer = append(pointsSpentPerPayer, PointsPerPayer{Payer: k, Points: v})
		transactionRecords = append(transactionRecords, Transaction{Payer: k, Points: v, TimeStamp: time.Now()})
	}

	context.IndentedJSON(http.StatusOK, pointsSpentPerPayer)
}

func getPointBalances(context *gin.Context) {
	pointBalances := calculatePointBalances()

	context.IndentedJSON(http.StatusOK, pointBalances)
}

// HELPER FUNCTIONS

func calculatePointBalances() map[string]int {
	pointBalances := make(map[string]int)

	for _, t := range transactionRecords {
		if _, exists := pointBalances[t.Payer]; exists {
			pointBalances[t.Payer]+=t.Points
		} else {
			pointBalances[t.Payer] = t.Points
		}
	}

	return pointBalances
}

func calculateNegativePointBalances() map[string]int {
	pointBalances := make(map[string]int)

	for _, t := range transactionRecords {
		if _, exists := pointBalances[t.Payer]; exists {
			if t.Points < 0 {
				pointBalances[t.Payer]+=t.Points
			}
		} else {
			if t.Points < 0 {
				pointBalances[t.Payer] = t.Points
			} else {
				pointBalances[t.Payer] = 0
			}
		}
	}

	return pointBalances
}

// MAIN FUNCTION

func main() {
	router := gin.Default()
	router.GET("/transaction", getTransactionRecords) // Route designed for testing purposes
	router.POST("/transaction", addTransactionRecord)
	router.PUT("/points", spendPoints)
	router.GET("/points", getPointBalances)
	router.Run("localhost:8080")
}
