package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

type Account interface {
	IsUserExist(ctx context.Context, username, password string) (*types.Account, error)
	CreateAccount(ctx context.Context, username, password string, isAdmin bool) (*types.Account, error)
	Seed(ctx context.Context) error
}

type account struct {
	accountRepo repos.Account
}

func NewAccount(accRepo repos.Account) Account {
	return &account{accountRepo: accRepo}
}

func (a *account) isUsernameTaken(ctx context.Context, username string) bool {
	_, err := a.accountRepo.FindByUsername(ctx, username)
	return err == nil
}

func (a *account) IsUserExist(ctx context.Context, username, password string) (*types.Account, error) {
	account, err := a.accountRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("wrong credentials")
	}
	if !utils.CheckPasswordHash(password, account.Password) {
		return nil, errors.New("wrong credentials")
	}
	return account, nil
}

func (a *account) CreateAccount(ctx context.Context, username, password string, isAdmin bool) (*types.Account, error) {
	pwdHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	if a.isUsernameTaken(ctx, username) {
		return nil, errors.New("username is already taken")
	}
	acc := types.Account{
		Id:       uuid.NewString(),
		Username: username,
		Password: pwdHash,
		Created:  time.Now().Format(time.RFC3339),
		IsAdmin:  isAdmin,
	}
	return a.accountRepo.Create(ctx, acc)
}

func (a *account) Seed(ctx context.Context) error {
	_, err := a.CreateAccount(ctx, "admin", "admin", true)
	return err
}
