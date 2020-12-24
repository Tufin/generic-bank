package common

import (
	"io"

	log "github.com/sirupsen/logrus"
)

func CloseWithErrLog(closer io.Closer) {

	if err := closer.Close(); err != nil {
		log.Errorf("failed to close with '%v'", err)
	}
}
