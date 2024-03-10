// cmd/api-server/seed/seed.go
package seed

import (
	"log"

	"github.com/shahin-bayat/go-api/internal/model"
	"github.com/shahin-bayat/go-api/internal/store"
)

// SeedAccounts seeds the database with initial accounts.
func SeedAccounts(store store.Store) {
	accounts := []*model.CreateAccountRequest{
		{
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password",
		},
		{
			FirstName: "Jane",
			LastName:  "Doe",
			Password:  "password",
		},
		{
			FirstName: "Alice",
			LastName:  "Doe",
			Password:  "password",
		},
	}

	for _, acc := range accounts {
		account, err := model.NewAccount(acc.FirstName, acc.LastName, acc.Password)
		if err != nil {
			log.Fatal(err)
		}

		if err := store.CreateAccount(account); err != nil {
			log.Fatal(err)
		}
	}
}
