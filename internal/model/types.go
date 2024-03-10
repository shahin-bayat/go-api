package model

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	IBAN              string    `json:"iban"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

type TransferRequest struct {
	Iban   string `json:"iban"`
	Amount int64  `json:"amount"`
}

func (a *Account) ValidatePassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw))
	return err == nil
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Iban     string `json:"iban"`
	Password string `json:"password"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	// go get golang.org/x/crypto/bcrypt
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encpw),
		IBAN:              generateIBAN(),
		Balance:           0,
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func generateIBAN() string {
	return fmt.Sprintf("DE%02d%04d%04d%04d%04d%04d", rand.Intn(100), rand.Intn(10000), rand.Intn(10000), rand.Intn(10000), rand.Intn(10000), rand.Intn(10000))
}
