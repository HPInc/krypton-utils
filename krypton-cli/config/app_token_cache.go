package config

import (
	"cli/common"
	"fmt"
)

const (
	AppTokenFile = "app_token"
)

var (
	// cache file path
	AppTokenCacheFile = fmt.Sprintf("%s/%s", getConfigDir(), AppTokenFile)
)

type AppTokenCache struct {
	tokenData []byte
}

func NewAppTokenCache() *AppTokenCache {
	return &AppTokenCache{
		tokenData: readCachedAppTokenData(),
	}
}

// set app token
func (a *AppTokenCache) Update(data string) error {
	if err := common.WriteFile(AppTokenCacheFile, []byte(data)); err != nil {
		return err
	}
	a.tokenData = []byte(data)
	return nil
}

// read cached data if any
func readCachedAppTokenData() []byte {
	bytes, err := common.ReadFile(AppTokenCacheFile)
	if err != nil {
		logger.Debugf("Error reading app token from cache: file=%s, err=%v",
			AppTokenCacheFile, err)
		return nil
	}
	return bytes
}

// get token data
func (a *AppTokenCache) GetAppToken() string {
	if a.tokenData != nil {
		return string(a.tokenData)
	}
	return ""
}
