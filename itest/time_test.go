package itest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/tufin/generic-bank/common"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {

	timeURLPrefix := AdminTimeURL()
	for _, currZone := range []string{"Asia/Jerusalem", "Asia/Tokyo", "Europe/London"} {
		timeURL := fmt.Sprintf("%s?zone=%s", timeURLPrefix, currZone)
		log.Infof("getting time for '%s'... '%s'", currZone, timeURL)
		response, err := http.Get(timeURL)
		require.NoError(t, err, timeURL)
		require.Equal(t, response.StatusCode, http.StatusOK)

		payload, err := ioutil.ReadAll(response.Body)
		require.NoError(t, err)
		var time map[string]string
		require.NoError(t, json.Unmarshal(payload, &time))
		log.Info(time)
		require.Equal(t, "OK", time["status"])

		common.CloseWithErrLog(response.Body)
	}
}
