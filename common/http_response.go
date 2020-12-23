package common

import (
	"encoding/json"
	"net/http"
	"reflect"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	HeaderAuthorization        = "Authorization"
	HeaderAccept               = "Accept"
	HeaderContentType          = "Content-Type"
	HeaderServer               = "Server"
	HeaderCookie               = "Cookie"
	ContentTypeApplicationJSON = "application/json"
	ContentTypeApplicationYAML = "application/yaml"
	ContentTypeApplicationXML  = "application/xml"
	ContentTypeApplicationForm = "application/x-www-form-urlencoded"
	ContentTypeTextPlain       = "text/plain"
)

// Notes (reference:https://golang.org/pkg/net/http/)
// Changing the header map after a call to WriteHeader (or
// Write) has no effect unless the modified headers are
// trailers.
// If WriteHeader has not yet been called, Write calls
// WriteHeader(http.StatusOK) before writing the data. If the Header
// does not contain a Content-Type line, Write adds a Content-Type set
// to the result of passing the initial 512 bytes of written data to
// DetectContentType. Additionally, if the total size of all written
// data is under a few KB and there are no Flush calls, the
// Content-Length header is added automatically.
func RespondWith(w http.ResponseWriter, r *http.Request, code int, response interface{}) {

	w.Header().Add("Cache-Control", "no-cache")
	if response != nil {
		accept := r.Header.Get(HeaderAccept)
		if accept == "application/pretty+json" {
			w.Header().Set(HeaderContentType, ContentTypeApplicationJSON)
			w.WriteHeader(code)
			if data, err := json.MarshalIndent(response, "", "\t"); err != nil {
				log.Errorf("failed to marshal indent response '%+v' with '%v'", response, err)
			} else {
				if _, err := w.Write(data); err != nil {
					log.Errorf("failed to write response '%+v' with '%v'", response, err)
				}
			}
		} else if accept == ContentTypeApplicationYAML {
			w.Header().Set(HeaderContentType, ContentTypeApplicationYAML)
			w.WriteHeader(code)
			if err := yaml.NewEncoder(w).Encode(response); err != nil {
				log.Errorf("failed to yaml encode response '%+v' with '%v'", response, err)
			}
		} else if accept == ContentTypeTextPlain {
			w.Header().Set(HeaderContentType, ContentTypeTextPlain)
			w.WriteHeader(code)
			kind := reflect.TypeOf(response).Kind()
			if kind == reflect.String {
				if _, err := w.Write([]byte(response.(string))); err != nil {
					log.Errorf("failed to write response '%+v' with '%v'", response, err)
				}
			} else {
				log.Errorf("failed to stream response with unexpected response kind '%v' for header '%s'='%s' (should be 'string')",
					kind, HeaderContentType, ContentTypeTextPlain)
			}
		} else {
			w.Header().Set(HeaderContentType, ContentTypeApplicationJSON)
			w.WriteHeader(code)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Errorf("failed to encode json response '%+v' with '%v'", response, err)
			}
		}
	} else {
		w.WriteHeader(code)
	}
}
