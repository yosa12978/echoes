package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/yosa12978/echoes/types"
)

var (
	store *sessions.CookieStore
	once  sync.Once
)

func SetupStore() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

// return session storage here
// or may be implement get/set functions for storage
func Get(r *http.Request, key string) (interface{}, error) {
	//change this once
	once.Do(SetupStore)

	session, err := store.Get(r, "user_store")
	return session.Values[key], err
}

func Set(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
	//change this to something better
	once.Do(SetupStore)

	session, err := store.Get(r, "user_store")
	if err != nil {
		return err
	}
	session.Values[key] = value
	return session.Save(r, w)
}

func GetInfo(r *http.Request) (*types.SessionInfo, error) {
	//change this to something better
	once.Do(SetupStore)

	session, err := store.Get(r, "user_store")
	if err != nil {
		return nil, err
	}
	userval := session.Values["account"]
	if userval == nil {
		return nil, errors.New("session is empty")
	}
	var info types.SessionInfo
	err = json.Unmarshal([]byte(userval.(string)), &info)
	return &info, err
}

func SetInfo(w http.ResponseWriter, r *http.Request, account *types.Account) error {
	//change this to something better
	once.Do(SetupStore)

	session, err := store.Get(r, "user_store")
	if err != nil {
		return err
	}

	if account == nil {
		session.Values["account"] = nil
		return session.Save(r, w)
	}

	var role = "USER"
	if account.IsAdmin {
		role = "ADMIN"
	}

	sessionInfo := types.SessionInfo{
		Username:  account.Username,
		Role:      role,
		Timestamp: time.Now().UnixNano(),
	}
	acc, err := json.Marshal(sessionInfo)
	if err != nil {
		return err
	}
	session.Values["account"] = string(acc)
	return session.Save(r, w)
}
