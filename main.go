package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
)

type QueryRestriction struct {
	Key        string `json:"key"`
	ValueRegex string `json:"valueRegex"`
}

type Config struct {
	URL               string             `json:"url"`
	QueryRestrictions []QueryRestriction `json:"queryRestrictions"`
}

func ParseConfigFromBase64(encodedConfig string) (config Config) {
	raw, _ := base64.StdEncoding.DecodeString(encodedConfig)
	json.Unmarshal(raw, &config)
	return
}

func (config Config) serverURL() *url.URL {
	uri, err := url.Parse(config.URL)
	if err != nil {
		log.Fatal("Sever URL is invalid :" + config.URL)
	}
	return uri
}

func ProxyHandler(config Config) http.HandlerFunc {
	log.Println(config)
	origin := httputil.NewSingleHostReverseProxy(config.serverURL())
	return func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		for _, restriction := range config.QueryRestrictions {
			value := request.Form.Get(restriction.Key)
			if !regexp.MustCompile(restriction.ValueRegex).MatchString(value) {
				http.Error(writer, "FORBIDDEN", 403)
				return
			}
		}
		origin.ServeHTTP(writer, request)
	}
}

func main() {
	handler := ProxyHandler(ParseConfigFromBase64(os.Getenv("CONFIG")))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handler))
}
