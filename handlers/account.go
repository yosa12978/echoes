package handlers

import (
	"context"
	"net/http"

	"github.com/yosa12978/echoes/repos"
	"github.com/yosa12978/echoes/session"
	"github.com/yosa12978/echoes/utils"
)

type Account interface {
	Login(ctx context.Context) http.Handler
	Signup(ctx context.Context) http.Handler
	Logout(ctx context.Context) http.Handler
}

type account struct {
	accountRepo repos.Account
}

func NewAccount(accountRepo repos.Account) Account {
	h := new(account)
	h.accountRepo = accountRepo
	return h
}

func (h *account) Login(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		account, err := h.accountRepo.FindByUsername(ctx, username)
		if err != nil {
			utils.RenderBlock(w, "alert", "user not found")
			return
		}
		if !utils.CheckPasswordHash(password, account.Password) {
			utils.RenderBlock(w, "alert", "user not found")
			return
		}
		if err := session.SetInfo(w, r, account); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		w.Header().Set("HX-Redirect", "/admin")
	})
}

func (h *account) Logout(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := session.GetInfo(r)
		if err != nil {
			http.Error(w, "you can't logout unless you logged in", http.StatusUnauthorized)
			return
		}
		session.SetInfo(w, r, nil)
		w.Header().Set("HX-Redirect", "/")
	})
}

func (h *account) Signup(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// i don't need this one at the time
	})
}
