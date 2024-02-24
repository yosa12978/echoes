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
	IsUserExist(ctx context.Context, username, password string) error
	CreateAccount(ctx context.Context, username, password string) (*types.Account, error)
}

type account struct {
	accountRepo repos.Account
}

func NewAccount(accRepo repos.Account) Account {
	return &account{accountRepo: accRepo}
}

func (a *account) IsUserExist(ctx context.Context, username, password string) error {
	account, err := a.accountRepo.FindByUsername(ctx, username)
	if err != nil {
		return errors.New("wrong credentials")
	}
	if !utils.CheckPasswordHash(password, account.Password) {
		return errors.New("wrong credentials")
	}
	return nil
}

func (a *account) CreateAccount(ctx context.Context, username, password string) (*types.Account, error) {
	pwdHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	acc := types.Account{
		Id:       uuid.NewString(),
		Username: username,
		Password: pwdHash,
		Created:  time.Now().Format(time.RFC3339),
		IsAdmin:  false,
	}
	return a.accountRepo.Create(ctx, acc)
}
