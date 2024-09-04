package endpoints

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/yosa12978/echoes/logging"
	"github.com/yosa12978/echoes/services"
	"github.com/yosa12978/echoes/session"
	"github.com/yosa12978/echoes/types"
	"github.com/yosa12978/echoes/utils"
)

func Login(logger logging.Logger, service services.Account) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		username := body["username"].(string)
		password := body["password"].(string)
		account, err := service.GetByCredentials(r.Context(), username, password)
		if err != nil {
			if errors.Is(err, types.ErrNotFound) {
				utils.RenderBlock(w, "alert", "wrong credentials")
				return
			}
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		if err := session.StartSession(r, w, *account); err != nil {
			utils.RenderBlock(w, "alert", err.Error())
			return
		}
		logger.Info("user %s logged in", username)
		w.Header().Set("HX-Redirect", "/admin")
	}
}
