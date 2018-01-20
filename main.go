package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"log"
	"regexp"
)

type QueryRestriction struct {
	Key        string
	ValueRegex string
}

type Config struct {
	URL               string
	QueryRestrictions []QueryRestriction
}

func (config Config) serverURL() *url.URL {
	u, err := url.Parse(config.URL)
	if err != nil {
		log.Fatal("Sever URL is invalid :" + config.URL)
	}
	return u
}

func ProxyHandler(config Config) http.HandlerFunc {
	server := httputil.NewSingleHostReverseProxy(config.serverURL())
	return func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		for _, restriction := range config.QueryRestrictions {
			value := request.Form.Get(restriction.Key)
			if !regexp.MustCompile(restriction.ValueRegex).MatchString(value) {
				http.Error(writer, "FORBIDDEN", 403)
				return
			}
		}
		server.ServeHTTP(writer, request)
	}
}

func main() {

}
