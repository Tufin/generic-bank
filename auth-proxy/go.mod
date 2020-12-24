module github.com/tufin/generic-bank/auth-proxy

go 1.15

require (
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20201205024021-ac21108117ac // indirect
	github.com/tufin/generic-bank/common v0.0.0-20201224122121-3ba94a82a8ff
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)

replace github.com/tufin/generic-bank/common => ../common
