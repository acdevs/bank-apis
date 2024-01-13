package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Account struct {
	ID        string `json:"id"`
	Number    int64  `json:"number"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Balance   int    `json:"balance"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateAccountResponse struct {
	ID string `json:"id"`
}

type Error struct {
	Message string `json:"message"`
}

func NewAccount(firstName, lastName string) *Account {
	year := time.Now().Year() % 100
	randomNumber := rand.Int63n(90000) + 10000

	accountID := fmt.Sprintf("%d%05d", year, randomNumber)

	return &Account{
		ID:        accountID,
		Number:    rand.Int63n(9000000000) + 1000000000,
		FirstName: firstName,
		LastName:  lastName,
		Balance:   0,
	}
}