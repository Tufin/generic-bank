module github.com/tufin/generic-bank/db

go 1.15

replace github.com/tufin/generic-bank/common => ../common

require (
	github.com/gorilla/mux v1.8.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.4.0
	github.com/tufin/generic-bank/common v0.0.0-00010101000000-000000000000
)
