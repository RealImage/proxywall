package main

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"reflect"
	"strconv"
)

func TestConfigParsing(t *testing.T) {

	tests := []struct {
		encoded        string
		expectedConfig Config
	}{
		{
			encoded: "ew0KICAidXJsIjogImh0dHBzOlwvXC90ZXN0Lmdvb2dsZS5jb20iLA0KICAicXVlcnlSZXN0cmljdGlvbnMiOiBbDQogICAgew0KICAgICAgImtleSI6ICJ1cmwiLA0KICAgICAgInZhbHVlUmVnZXgiOiAiXFwuZ29vZ2xlXFwuY29tJCINCiAgICB9DQogIF0NCn0=",
			expectedConfig: Config{
				URL: "https://test.google.com",
				QueryRestrictions: []QueryRestriction{
					{
						Key:        "url",
						ValueRegex: "\\.google\\.com$",
					},
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run("Encoding "+strconv.Itoa(i+1), func(t *testing.T) {
			if !reflect.DeepEqual(tt.expectedConfig, ParseConfigFromBase64(tt.encoded)) {
				t.Error("Config parse error, expected: ", tt.expectedConfig, " got: ", ParseConfigFromBase64(tt.encoded))
			}
		})
	}
}

func TestProxyHandler(t *testing.T) {
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	proxyHandler := ProxyHandler(Config{
		URL: origin.URL,
		QueryRestrictions: []QueryRestriction{
			{
				Key:        "url",
				ValueRegex: "\\.google\\.com$",
			},
		},
	})
	proxy := httptest.NewServer(http.HandlerFunc(proxyHandler))

	type fields struct {
		query            string
		originStatusCode int
		proxyStatusCode  int
		originBody       string
		proxyBody        string
	}

	tests := []fields{
		{
			query:            "?url=test.google.com",
			originStatusCode: 200,
			proxyStatusCode:  200,
			originBody:       "OK",
			proxyBody:        "OK",
		},
		{
			query:            "?url=test.somethingelse.com",
			originStatusCode: 200,
			proxyStatusCode:  403,
			originBody:       "OK",
			proxyBody:        "FORBIDDEN",
		},
	}
	for _, exp := range tests {
		originResponse, _ := http.Get(origin.URL + exp.query)
		if originResponse.StatusCode != exp.originStatusCode {
			t.Fatal("Origin status code did not match", )
		}
		proxyResponse, _ := http.Get(proxy.URL + exp.query)
		if proxyResponse.StatusCode != exp.proxyStatusCode {
			t.Fatal("proxy status code did not match")
		}
	}
}
