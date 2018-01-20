package main

import (
	"net/http"
	"testing"
	"net/http/httptest"
)

//func TestConfig_serverURL(t *testing.T) {
//	type fields struct {
//		URL               string
//		QueryRestrictions []QueryRestriction
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   *url.URL
//	}{
//	// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			config := Config{
//				URL:               tt.fields.URL,
//				QueryRestrictions: tt.fields.QueryRestrictions,
//			}
//			if got := config.serverURL(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Config.serverURL() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

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
		query string
		originStatusCode int
		proxyStatusCode int
		originBody string
		proxyBody string
	}

	tests := []fields{
		{
			query: "?url=test.google.com",
			originStatusCode: 200,
			proxyStatusCode: 200,
			originBody: "OK",
			proxyBody: "OK",
		},
		{
			query: "?url=test.somethingelse.com",
			originStatusCode: 200,
			proxyStatusCode: 403,
			originBody: "OK",
			proxyBody: "FORBIDDEN",
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

//func Test_main(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//	// TODO: Add test cases.
//	}
//	for range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			main()
//		})
//	}
//}
