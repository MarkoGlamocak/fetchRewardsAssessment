package helpers

import (
	"github.com/MarkoGlamocak/fetchRewardsAssessment/models"
)

func CalculatePointBalances() map[string]int {
	pointBalances := make(map[string]int)

	for _, t := range models.TransactionRecords {
		if _, exists := pointBalances[t.Payer]; exists {
			pointBalances[t.Payer]+=t.Points
		} else {
			pointBalances[t.Payer] = t.Points
		}
	}

	return pointBalances
}
