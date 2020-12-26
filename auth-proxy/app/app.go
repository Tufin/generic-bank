package app

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
	"github.com/tufin/generic-bank/common"
)

var (
	Store        *sessions.FilesystemStore
	Domain       string
	ClientID     string
	ClientSecret string
	CallbackURL  string
)

func Init() {

	Store = sessions.NewFilesystemStore("", []byte("something-very-secret"))
	gob.Register(map[string]interface{}{})

	Domain = common.GetEnvOrExit("AUTH0_DOMAIN")
	ClientID = common.GetEnvSensitiveOrExit("AUTH0_CLIENT_ID")
	ClientSecret = common.GetEnvSensitiveOrExit("AUTH0_CLIENT_SECRET")
	CallbackURL = common.GetEnvOrExit("CALLBACK_URL")
}
