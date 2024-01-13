package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)


type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	GetAccountByID(string) (*Account, error)
	DeleteAccount(string) error
}


type PostgresStorage struct {
	db *sql.DB
}


func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=2404 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}


func (s *PostgresStorage) Init() error{
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id VARCHAR(10) PRIMARY KEY,
		number BIGINT UNIQUE NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		balance BIGINT NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`
	_, err := s.db.Exec(query)
	return err
}


func (s *PostgresStorage) CreateAccount(account *Account) error {
	query := `INSERT INTO accounts (id, number, first_name, last_name, balance) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(query, account.ID, account.Number, account.FirstName, account.LastName, account.Balance)
	return err
}


func (s *PostgresStorage) GetAccountByID(id string) (*Account, error) {
	query := `SELECT * FROM accounts WHERE id = $1`
	row := s.db.QueryRow(query, id)

	account := new(Account)
	err := row.Scan(&account.ID, &account.Number, &account.FirstName, &account.LastName, &account.Balance)

	return account, err
}


func (s *PostgresStorage) UpdateAccount(account *Account) error {
	// _, err := s.db.Exec()
	return nil
}


func (s *PostgresStorage) DeleteAccount(id string) error {
	// _, err := s.db.Exec()
	return nil
}