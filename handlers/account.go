package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/session"
	"github.com/yosa12978/echoes/utils"
)

type Account interface {
	Login(ctx context.Context) http.Handler
	Signup(ctx context.Context) http.Handler
	Logout(ctx context.Context) http.Handler
}

type account struct {
	accountService services.Account
	logger         logging.Logger
}

func NewAccount(accountService services.Account, logger logging.Logger) Account {
	h := new(account)
	h.accountService = accountService
	h.logger = logger
	return h
}

func (h *account) Login(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}

		username := body["username"].(string)
		password := body["password"].(string)

		// r.ParseForm()
		// username := r.FormValue("username")
		// password := r.FormValue("password")
		account, err := h.accountService.IsUserExist(ctx, username, password)
		if err != nil {
			utils.RenderBlock(w, "alert", "user not found")
			return
		}
		if err := session.SetInfo(w, r, account); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		h.logger.Printf("user %s logged in", username)
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
