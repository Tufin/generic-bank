package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"github.com/tufin/generic-bank/common"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/admin/accounts", func(w http.ResponseWriter, r *http.Request) {
		log.Info("/admin/accounts")
		if _, err := w.Write([]byte("accounts :)")); err != nil {
			log.Errorf("failed to stream response with '%v'", err)
		}
	}).Methods(http.MethodGet)
	router.HandleFunc("/customer/balance", func(w http.ResponseWriter, r *http.Request) {
		log.Info("GET /customer/balance")
		if _, err := w.Write([]byte("123")); err != nil {
			log.Errorf("failed to stream response with '%v'", err)
		}
	})

	common.ServeMulti([]*mux.Router{router}, []string{"3000"})
}
