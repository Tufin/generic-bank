package auth

import (
	"encoding/json"
	"net/http"

	"github.com/tufin/generic-bank/auth-proxy/app"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(session.Values["profile"])
}
