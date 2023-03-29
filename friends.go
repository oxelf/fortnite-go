package main

import "time"

type Friend struct {
	AccountId string    `json:"accountId"`
	Status    string    `json:"status"`
	Direction string    `json:"direction"`
	Created   time.Time `json:"created"`
	Favorite  bool      `json:"favorite"`
}
