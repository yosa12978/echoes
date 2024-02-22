package services

import "github.com/yosa12978/echoes/repos"

type Account interface {
}

type account struct {
	accountRepo repos.Account
}

func NewAccount(accRepo repos.Account) Account {
	return &account{accountRepo: accRepo}
}
