package auth

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tufin/generic-bank/auth-proxy/app"
)

func IsAuthenticate(r *http.Request) bool {

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		log.Errorf("failed to get auth session with '%v'", err)
		return false
	}
	if _, ok := session.Values["profile"]; !ok {
		//http.Redirect(w, r, "/", http.StatusSeeOther)
		return false
	}

	return true
}
