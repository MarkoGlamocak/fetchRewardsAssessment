package helpers

import (
	"github.com/MarkoGlamocak/fetchRewardsAssessment/models"
)

func CalculateNegativePointBalances() map[string]int {
	pointBalances := make(map[string]int)

	for _, t := range models.TransactionRecords {
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
