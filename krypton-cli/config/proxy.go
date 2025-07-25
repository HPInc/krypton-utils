package config

import (
	"os"
)

const (
	EnvHttpProxy  = "HTTP_PROXY"
	EnvHttpsProxy = "HTTPS_PROXY"
)

// set proxy env
func (p *Proxy) apply() {
	_ = os.Setenv(EnvHttpProxy, p.HttpProxy)
	_ = os.Setenv(EnvHttpsProxy, p.HttpsProxy)
}
