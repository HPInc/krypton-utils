package config

import (
	"cli/common"
	"encoding/json"
	"fmt"
)

const (
	BulkEnrollFile = "bulk_enroll"
)

var (
	// cache file path
	BulkEnrollCacheFile = fmt.Sprintf("%s/%s", getConfigDir(), BulkEnrollFile)
)

// cached token
type bulkEnrollData struct {
	TenantId string `json:"tenant_id"`
	Token    string `json:"enroll_token"`
}

type BulkEnrollCache struct {
	tokenData *bulkEnrollData
}

func NewBulkEnrollCache() *BulkEnrollCache {
	return &BulkEnrollCache{
		tokenData: readCachedBulkEnrollData(),
	}
}

// set bulk enroll cache
func (e *BulkEnrollCache) Update(bytes []byte) error {
	if err := common.WriteFile(BulkEnrollCacheFile, bytes); err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, e.tokenData); err != nil {
		return err
	}
	return nil
}

// read cached data if any
func readCachedBulkEnrollData() *bulkEnrollData {
	et := bulkEnrollData{}
	bytes, err := common.ReadFile(BulkEnrollCacheFile)
	if err != nil {
		logger.Debugf("Error reading bulk enroll token from cache: file=%s, err=%v",
			BulkEnrollCacheFile, err)
		return nil
	}
	if err = json.Unmarshal(bytes, &et); err != nil {
		logger.Debugf("Error unmarshal bulk enroll token from cache: file=%s, err=%v",
			BulkEnrollCacheFile, err)
		return nil
	}
	return &et
}

// get token data
func (e *BulkEnrollCache) GetToken() string {
	if e.tokenData != nil {
		return e.tokenData.Token
	}
	return ""
}

// get tid claim from bulk enroll file
func (e *BulkEnrollCache) GetTenantId() string {
	if e.tokenData == nil || e.tokenData.TenantId == "" {
		return e.tokenData.TenantId
	}
	return ""
}
