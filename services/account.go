package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yosa12978/echoes/config"
	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

type Account interface {
	GetByCredentials(ctx context.Context, username, password string) (*types.Account, error)
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

func (a *account) GetByCredentials(ctx context.Context, username, password string) (*types.Account, error) {
	account, err := a.accountRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("wrong credentials")
	}
	fmt.Printf("\"%s\"\n", account.Password)
	if !utils.CheckPasswordHash(password+account.Salt, account.Password) {
		return nil, errors.New("wrong credentials")
	}
	return account, nil
}

func (a *account) CreateAccount(ctx context.Context, username, password string, isAdmin bool) (*types.Account, error) {
	if a.isUsernameTaken(ctx, username) {
		return nil, errors.New("username is already taken")
	}
	salt := uuid.NewString()
	pwdHash, err := utils.HashPassword(password + salt)
	if err != nil {
		return nil, err
	}

	acc := types.Account{
		Id:       uuid.NewString(),
		Username: username,
		Password: pwdHash,
		Created:  time.Now().UTC().Format(time.RFC3339),
		IsAdmin:  isAdmin,
		Salt:     salt,
	}
	return a.accountRepo.Create(ctx, acc)
}

// refactor this
func (a *account) Seed(ctx context.Context) error {
	cfg := config.Get()
	usr, err := a.accountRepo.FindByUsername(ctx, "root")
	if err != nil {
		if _, err := a.CreateAccount(
			ctx,
			"root",
			cfg.Server.RootPass,
			true,
		); err != nil {
			return err
		}

	}
	if usr != nil {
		if !utils.CheckPasswordHash(cfg.Server.RootPass+usr.Salt, usr.Password) {
			a.changeRootPassword(ctx, usr.Id, cfg.Server.RootPass)
		}
	}
	return nil
}

// remove userId
func (a *account) changeRootPassword(ctx context.Context, userId, newPassword string) error {
	fmt.Println("changing password")
	if strings.Contains(newPassword, " ") {
		return errors.New("password can't contain spaces")
	}
	if len(newPassword) < 4 {
		return errors.New("length of your password can't be less then 4 characters")
	}
	salt := uuid.NewString()
	passwordHash, _ := utils.HashPassword(newPassword + salt)
	return a.accountRepo.Update(ctx, userId, types.Account{
		Username: "root",
		IsAdmin:  true,
		Password: passwordHash,
		Salt:     salt,
	})
}
