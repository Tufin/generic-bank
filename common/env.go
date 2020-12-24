package common

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func GetEnvSensitive(variable string) string {

	ret := os.Getenv(variable)
	if ret != "" {
		log.Infof("'%s': [sensitive]", variable)
	}

	return ret
}

func GetEnvSensitiveOrExit(variable string) string {

	ret := GetEnvSensitive(variable)
	if ret == "" {
		log.Fatalf("Please, set '%s'", variable)
	}

	return ret
}

func GetEnvOrExit(variable string) string {

	ret := os.Getenv(variable)
	if ret == "" {
		log.Fatalf("Please, set '%s'", variable)
	}
	log.Infof("'%s': '%s'", variable, ret)

	return ret
}

func GetEnvWithDefault(variable, defaultValue string) string {

	ret := os.Getenv(variable)
	if ret == "" {
		ret = defaultValue
	}
	log.Infof("'%s': '%s'", variable, ret)

	return ret
}
