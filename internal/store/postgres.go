package store

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/shahin-bayat/go-api/internal/model"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
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

	return &PostgresStore{db: db}, nil

}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
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

func (s *PostgresStore) CreateAccount(a *model.Account) error {
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

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM account WHERE id = $1", id)
	return err
}

func (s *PostgresStore) UpdateAccount(id int, a *model.Account) (*model.Account, error) {
	query := `
	UPDATE account 
	SET first_name = $1, last_name = $2, iban = $3, encrypted_password = $4, balance = $5, created_at = $6
	WHERE id = $7
	`
	_, err := s.db.Exec(query, a.FirstName, a.LastName, a.IBAN, a.EncryptedPassword, a.Balance, a.CreatedAt, id)
	if err != nil {
		return nil, err
	}
	return s.GetAccountById(id)
}

func (s *PostgresStore) GetAccounts() ([]*model.Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}
	accounts := []*model.Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *PostgresStore) GetAccountByIban(iban string) (*model.Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE iban = $1", iban)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account with iban %s not found", iban)
}

func (s *PostgresStore) GetAccountById(id int) (*model.Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*model.Account, error) {
	account := new(model.Account)
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
