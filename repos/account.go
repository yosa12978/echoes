package repos

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/types"
)

type Account interface {
	FindById(ctx context.Context, id string) (*types.Account, error)
	FindByCredentials(ctx context.Context, username, passwordHash string) (*types.Account, error)
	FindByUsername(ctx context.Context, username string) (*types.Account, error)
	Create(ctx context.Context, account types.Account) (*types.Account, error)
	Update(ctx context.Context, accountId string, account types.Account) error
	Delete(ctx context.Context, accountId string) error
}

type account struct {
	db *sql.DB
}

func NewAccountPostgres() Account {
	repo := new(account)
	repo.db = data.Postgres()
	return repo
}

// type Account struct {
// 	Id       string
// 	Username string
// 	Password string
// 	Created  string
//  IsAdmin bool
// }

func (repo *account) FindById(ctx context.Context, id string) (*types.Account, error) {
	var acc types.Account
	q := "SELECT * FROM accounts WHERE id=$1;"
	err := repo.db.QueryRowContext(ctx, q, id).Scan(&acc.Id, &acc.Username, &acc.Password, &acc.Created, &acc.IsAdmin, &acc.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &acc, nil
}

// this works wrong
func (repo *account) FindByCredentials(ctx context.Context, username, passwordHash string) (*types.Account, error) {
	var acc types.Account
	q := "SELECT * FROM accounts WHERE username=$1 AND password=$2;"
	err := repo.db.
		QueryRowContext(ctx, q, strings.ToLower(username), passwordHash).
		Scan(&acc.Id, &acc.Username, &acc.Password, &acc.Created, &acc.IsAdmin, &acc.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &acc, nil
}

func (repo *account) FindByUsername(ctx context.Context, username string) (*types.Account, error) {
	var acc types.Account
	q := "SELECT * FROM accounts WHERE username=$1;"
	err := repo.db.
		QueryRowContext(ctx, q, strings.ToLower(username)).
		Scan(&acc.Id, &acc.Username, &acc.Password, &acc.Created, &acc.IsAdmin, &acc.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return &acc, nil
}

func (repo *account) Create(ctx context.Context, account types.Account) (*types.Account, error) {
	q := "INSERT INTO accounts (id, username, password, created, isadmin, salt) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err := repo.db.ExecContext(ctx, q,
		account.Id,
		strings.ToLower(account.Username),
		account.Password, account.Created,
		account.IsAdmin,
		account.Salt)
	if err != nil {
		return nil, errors.Join(err, ErrInternalFailure)
	}
	return nil, nil
}

func (repo *account) Update(ctx context.Context, accountId string, account types.Account) error {
	q := "UPDATE accounts SET username=$1, password=$2, isadmin=$3, salt=$4 WHERE id=$5;"
	_, err := repo.db.ExecContext(ctx, q,
		strings.ToLower(account.Username),
		account.Password,
		account.IsAdmin,
		account.Salt,
		accountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return errors.Join(err, ErrInternalFailure)
	}
	return nil
}

func (repo *account) Delete(ctx context.Context, accountId string) error {
	q := "DELETE FROM accounts WHERE id=$1;"
	_, err := repo.db.ExecContext(ctx, q, accountId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return errors.Join(err, ErrInternalFailure)
	}
	return nil
}
