package models

import "time"

type Transaction struct {
	Payer string `json:"payer"`
	Points int `json:"points"`
	TimeStamp time.Time `json:"timestamp"`
}
