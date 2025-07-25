package common

import (
	"fmt"
	"net/http"
)

const (
	Authorization = "Authorization"
	UserAgent     = "User-Agent"

	KryptonUserAgent = "krypton-cli"
)

// generic add header
func AddRequestHeader(req *http.Request, key, value string) {
	req.Header.Add(key, value)
}

// add user-agent to http request
func AddUserAgentHeader(req *http.Request) {
	req.Header.Add(UserAgent, KryptonUserAgent)
}

// add authorization header to http request
func AddAuthorizationHeader(req *http.Request, bearerToken string) {
	if bearerToken != "" {
		req.Header.Add(Authorization,
			fmt.Sprintf("Bearer %s", bearerToken))
	}
}
