package session

import (
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/yosa12978/echoes/configs"
	"github.com/yosa12978/echoes/types"
)

var (
	store *sessions.CookieStore
)

func init() {
	gob.Register(types.Session{})
}

func SetupStore() {
	store = sessions.NewCookieStore([]byte(configs.Get().SessionKey))
}

func Get(r *http.Request, key string) (interface{}, error) {
	session, err := store.Get(r, "echoes_session")
	if err != nil {
		return nil, err
	}
	return session.Values[key], nil
}

func Set(r *http.Request, w http.ResponseWriter, key string, value interface{}) error {
	session, err := store.Get(r, "echoes_session")
	if err != nil {
		return err
	}
	session.Values[key] = value
	return session.Save(r, w)
}

func Delete(r *http.Request, w http.ResponseWriter, key string) error {
	session, err := store.Get(r, "echoes_session")
	if err != nil {
		return err
	}
	delete(session.Values, key)
	return session.Save(r, w)
}

func EndSession(r *http.Request, w http.ResponseWriter) error {
	session, err := store.Get(r, "echoes_session")
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func StartSession(r *http.Request, w http.ResponseWriter, account types.Account) error {
	session, err := store.New(r, "echoes_session")
	if err != nil {
		return err
	}
	session.Values["account"] = types.Session{
		Username:        account.Username,
		IsAdmin:         account.IsAdmin,
		IsAuthenticated: true,
		Timestamp:       time.Now().UTC().UnixNano(),
	}
	return session.Save(r, w)
}

func GetSession(r *http.Request) (*types.Session, error) {
	session, err := store.Get(r, "echoes_session")
	if err != nil {
		return nil, err
	}
	if value, ok := session.Values["account"].(types.Session); ok {
		return &value, nil
	}
	return nil, errors.New("user is not logged in")
}
