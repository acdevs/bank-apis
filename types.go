package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID        string `json:"id"`
	Number    int64  `json:"number"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Balance   int64    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct  {
	ID string `json:"id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password string `json:"password"`
}

type CreateAccountResponse struct {
	ID string `json:"id"`
}

type DeleteAccountResponse struct {
	ID string `json:"id"`
}

type UpdateAccountBalanceRequest struct {
	Amount int64 `json:"amount"`
}

type TransferAccountRequest struct {
	ToAccount   int64 `json:"to_account_number"`
	Amount      int64    `json:"amount"`
}

type Response struct {
	Message string `json:"message"`
}

type Error struct {
	Message string `json:"error"`
}

func NewAccount(firstName, lastName, password string) *Account {
	year := time.Now().Year() % 100
	randomNumber := rand.Int63n(90000) + 10000

	accountID := fmt.Sprintf("%d%05d", year, randomNumber)

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating password", err)
	}

	return &Account{
		ID:        accountID,
		Number:    rand.Int63n(9000000000) + 1000000000,
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(encryptedPassword),
		Balance:   0,
		CreatedAt: time.Now(),
	}
}