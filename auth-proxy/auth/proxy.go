package auth

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/tufin/generic-bank/common"
)

func CreateAuthProxy() *httputil.ReverseProxy {

	target := common.GetEnvOrExit("TARGET_URL")
	origin, err := url.Parse(target)
	if err != nil {
		log.Fatalf("failed to parse proxy target URL '%s' with '%v'", target, err)
	}

	return &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host
	}}
}

//func CreateAuthProxy(proxyPath string) *httputil.ReverseProxy {
//
//	const proxyURL = "http://localhost:8081"
//	origin, err := url.Parse(proxyURL)
//	if err != nil {
//		log.Fatalf("failed to parse proxy URL '%s' with '%v'", proxyURL, err)
//	}
//	ret := httputil.NewSingleHostReverseProxy(origin)
//	ret.Director = func(req *http.Request) {
//		req.Header.Add("X-Forwarded-Host", req.Host)
//		req.Header.Add("X-Origin-Host", origin.Host)
//		req.URL.Scheme = origin.Scheme
//		req.URL.Host = origin.Host
//
//		wildcardIndex := strings.IndexAny(proxyPath, "*")
//		proxyPath := singleJoiningSlash(origin.Path, req.URL.Path[wildcardIndex:])
//		if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
//			proxyPath = proxyPath[:len(proxyPath)-1]
//		}
//		req.URL.Path = origin.Path
//	}
//
//	return ret
//}

//func singleJoiningSlash(a, b string) string {
//
//	aslash := strings.HasSuffix(a, "/")
//	bslash := strings.HasPrefix(b, "/")
//	switch {
//	case aslash && bslash:
//		return a + b[1:]
//	case !aslash && !bslash:
//		return a + "/" + b
//	}
//
//	return a + b
//}
//
//func CreateAuthProxy(target string) *httputil.ReverseProxy {
//
//	targetURL, err := url.Parse(target)
//	if err != nil {
//		log.Fatalf("failed to parse reverse proxy url '%s' with '%v'", target, err)
//	}
//
//	targetQuery := targetURL.RawQuery
//	director := func(req *http.Request) {
//		req.URL.Scheme = targetURL.Scheme
//		req.URL.Host = targetURL.Host
//		req.Host = targetURL.Host
//		req.URL.Path = singleJoiningSlash(targetURL.Path, req.URL.Path)
//		if targetQuery == "" || req.URL.RawQuery == "" {
//			req.URL.RawQuery = targetQuery + req.URL.RawQuery
//		} else {
//			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
//		}
//		if _, ok := req.Header["User-Agent"]; !ok {
//			// explicitly disable User-Agent so it's not set to default value
//			req.Header.Set("User-Agent", "")
//		}
//	}
//
//	return &httputil.ReverseProxy{Director: director}
//}
