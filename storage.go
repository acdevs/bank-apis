package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)


type Storage interface {
	CreateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(string) (*Account, error)
	DeleteAccount(string) error
	UpdateAccountBalance(string, int64) error
	TransferAccount(int64, int64, int64) error
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
		password VARCHAR(255) NOT NULL,
		balance BIGINT NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`
	_, err := s.db.Exec(query)
	return err
}


func (s *PostgresStorage) CreateAccount(account *Account) error {
	query := `INSERT INTO accounts (id, number, first_name, last_name, password, balance, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.Exec(query, account.ID, account.Number, account.FirstName, account.LastName, account.Password, account.Balance, account.CreatedAt)
	return err
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM accounts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, 0)
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(&account.ID, &account.Number, &account.FirstName, &account.LastName, &account.Password, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStorage) GetAccountByID(id string) (*Account, error) {
	query := `SELECT * FROM accounts WHERE id = $1`
	row := s.db.QueryRow(query, id)

	account := new(Account)
	err := row.Scan(&account.ID, &account.Number, &account.FirstName, &account.LastName, &account.Password, &account.Balance, &account.CreatedAt)

	return account, err
}


func (s *PostgresStorage) DeleteAccount(id string) error {
	query := "DELETE FROM accounts WHERE id = $1"
	_, err := s.db.Exec(query, id)
	return err
}


func (s *PostgresStorage) UpdateAccountBalance(id string, balance int64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err := s.db.Exec(query, balance, id)
	return err
}


func (s *PostgresStorage) TransferAccount(fromNumber, toNumber, amount int64) error {
	query := `UPDATE accounts SET balance = balance - $1 WHERE number = $2`
	_, err := s.db.Exec(query, amount, fromNumber)
	if err != nil {
		return err
	}

	query = `UPDATE accounts SET balance = balance + $1 WHERE number = $2`
	_, err = s.db.Exec(query, amount, toNumber)
	return err
}