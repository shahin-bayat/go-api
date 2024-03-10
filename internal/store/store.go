package store

import "github.com/shahin-bayat/go-api/internal/model"

type Store interface {
	CreateAccount(*model.Account) error
	DeleteAccount(int) error
	UpdateAccount(*model.Account) error
	GetAccounts() ([]*model.Account, error)
	GetAccountById(int) (*model.Account, error)
	GetAccountByIban(string) (*model.Account, error)
}
