module github.com/tufin/generic-bank

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/prometheus/common v0.15.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/tufin/generic-bank/common v0.0.0-00010101000000-000000000000
)

replace github.com/tufin/generic-bank/common => ./common
