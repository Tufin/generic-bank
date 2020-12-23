package common

import (
	"os"

	log "github.com/sirupsen/logrus"
	sabik "github.com/tufin/sabik/client"
	"github.com/tufin/sabik/common/env"
)

func CreateMiddleware(serviceName string) *sabik.Middleware {

	fatalOnError(os.Setenv(env.KeyDomain, "generic-bank"))
	fatalOnError(os.Setenv(env.KeyProject, "retail"))
	if serviceName != "" {
		fatalOnError(os.Setenv(sabik.EnvKeyServiceName, serviceName))
	}

	return sabik.NewMiddleware()
}

func fatalOnError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}
