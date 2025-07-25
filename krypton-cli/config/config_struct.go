package config

// auth related config
type Auth struct {
	TokenType        string `yaml:"token_type"`
	ManagementServer string `yaml:"mgmt_server"`
}

// proxy config
type Proxy struct {
	HttpProxy  string `yaml:"http_proxy"`
	HttpsProxy string `yaml:"https_proxy"`
}

// server addresses for destination servers
type Server struct {
	Addresses map[string]string `yaml:"addresses"`
}

// profile bag holding all configs
type Profile struct {
	Name   string `yaml:"name"`
	Auth   Auth   `yaml:"auth"`
	Proxy  Proxy  `yaml:"proxy"`
	Server Server `yaml:"server"`
}

// test strings for file upload tests
type EncryptedTestData struct {
	Data string `yaml:"data"` // encrypted data
	Key  string `yaml:"key"`  // encryption key
}

// Config holds all profiles and overrides to current profile
type Config struct {
	Auth              string    `yaml:"auth"`
	ProfileName       string    `yaml:"default_profile"`
	Profiles          []Profile `yaml:"profiles"`
	CurrentProfile    *Profile
	EnrollTokenCache  *EnrollTokenCache
	DeviceTokenCache  *DeviceTokenCache
	AppTokenCache     *AppTokenCache
	BulkEnrollCache   *BulkEnrollCache
	EncryptedTestData map[string]EncryptedTestData `yaml:"encrypted_test_data"`
}
