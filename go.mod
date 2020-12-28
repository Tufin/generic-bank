module github.com/tufin/generic-bank

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/prometheus/common v0.15.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/tufin/generic-bank/common v0.0.0-20201224122121-3ba94a82a8ff
	github.com/tufin/sabik v0.0.0-20201228104709-688081faec97
)

replace github.com/tufin/generic-bank/common => ./common
