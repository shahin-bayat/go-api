package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	store, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding the db...")
		seedAccounts(store)
	}
	server := NewAPIServer(":8000", store)
	server.Run()
}

func seedAccounts(store Storage) {
	accounts := []*CreateAccountRequest{
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

		account, err := NewAccount(acc.FirstName, acc.LastName, acc.Password)
		if err != nil {
			log.Fatal(err)
		}

		if err := store.CreateAccount(account); err != nil {
			log.Fatal(err)
		}
	}
}
