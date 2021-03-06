package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/tufin/generic-bank/auth-proxy/app"
	"golang.org/x/oauth2"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state, oauth2.SetAuthURLParam("audience", "http://localhost:8080/admin/accounts")), http.StatusTemporaryRedirect)
}
