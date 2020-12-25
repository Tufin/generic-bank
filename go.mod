module github.com/tufin/generic-bank

go 1.15

require (
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20201205024021-ac21108117ac // indirect
	github.com/prometheus/common v0.15.0
	github.com/sirupsen/logrus v1.7.0
	github.com/tufin/generic-bank/common v0.0.0-00010101000000-000000000000
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)

replace github.com/tufin/generic-bank/common => ./common
