package downstream

import (
	"crypto/tls"
	"net/http"
)

var (
	suggestionEngineClient *http.Client
)

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{ServerName: "suggestions.skc-ygo-api.com"},
	}
	suggestionEngineClient = &http.Client{Transport: tr}
}
