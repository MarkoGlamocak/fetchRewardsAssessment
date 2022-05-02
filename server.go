package main

import (
	"github.com/MarkoGlamocak/fetchRewardsAssessment/routers"
)

func main() {
	router := routers.Init()
	router.Run("localhost:8080")
}
