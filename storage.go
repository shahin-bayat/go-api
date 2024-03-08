package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByIban(string) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	// https://pkg.go.dev/github.com/lib/pq
	// from docker: docker run --name postgres -e POSTGRES_PASSWORD=goapi -p 5432:5432 -d postgres
	connStr := "user=postgres dbname=postgres password=goapi sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &PostgresStorage{db: db}, nil

}

func (s *PostgresStorage) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStorage) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		iban VARCHAR(50) UNIQUE,
		encrypted_password VARCHAR(100),
		balance INTEGER,
		created_at TIMESTAMP DEFAULT current_timestamp
)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateAccount(a *Account) error {
	query := `
	INSERT INTO account 
	(first_name, last_name, iban, encrypted_password, balance, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.Exec(query, a.FirstName, a.LastName, a.IBAN, a.EncryptedPassword, a.Balance, a.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM account WHERE id = $1", id)
	return err
}

func (s *PostgresStorage) UpdateAccount(a *Account) error {
	return nil
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *PostgresStorage) GetAccountByIban(iban string) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE iban = $1", iban)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account with iban %s not found", iban)
}

func (s *PostgresStorage) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.IBAN,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt)
	return account, err

}
