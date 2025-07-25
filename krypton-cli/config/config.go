package config

import (
	"cli/common"
	"cli/logging"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	KryptonConfigDir = ".krypton-cli"
	EnvConfigFile    = "CLI_CONFIG_FILE"
	IotRootCertFile  = "aws_root.cert"
	ConfigFileName   = "config.yaml"
)

var (
	DefaultConfigFile = filepath.Join(getConfigDir(), ConfigFileName)
	settings          Config
	logger            *logging.Log
)

func init() {
	logger = logging.InitLogger(logging.Info)
	if !Load() {
		os.Exit(1)
	}
}

// load config
func Load() bool {
	configFile := getConfigFile()

	logger.Debug("loading config file: ", configFile)

	// Open the configuration file for parsing.
	bytes, err := common.ReadFile(configFile)
	if err != nil {
		logger.Println("Failed to load configuration file: ",
			configFile, ", Error = ", err)
		return false
	}

	// Read the configuration file and unmarshal the YAML.
	err = yaml.Unmarshal(bytes, &settings)
	if err != nil {
		logger.Println("Failed to parse configuration file: ",
			configFile, ", Error = ", err)
		return false
	}

	// override config from environment variables
	// note this only happens if environment variables are specified
	settings.OverrideFromEnvironment()

	// apply current profile
	settings.applyProfile()
	//logger.Println(settings)

	settings.EnrollTokenCache = NewEnrollTokenCache()
	settings.DeviceTokenCache = NewDeviceTokenCache()
	settings.AppTokenCache = NewAppTokenCache()
	settings.BulkEnrollCache = NewBulkEnrollCache()

	return settings.CurrentProfile != nil
}

// read config file from cmdline args
// note: this will only work as an arg to the krypton-cli program
// not to the modules
func getConfigFile() string {
	configFile := os.Getenv(EnvConfigFile)
	if configFile == "" {
		configFile = DefaultConfigFile
	}
	return configFile
}

// return settings loaded from config
func GetSettings() *Config {
	return &settings
}

// return enroll token cache handler
func GetEnrollTokenCache() *EnrollTokenCache {
	return settings.EnrollTokenCache
}

// return device token cache handler
func GetDeviceTokenCache() *DeviceTokenCache {
	return settings.DeviceTokenCache
}

// return app token cache handler
func GetAppTokenCache() *AppTokenCache {
	return settings.AppTokenCache
}

// return bulk enroll cache handler
func GetBulkEnrollCache() *BulkEnrollCache {
	return settings.BulkEnrollCache
}

// return config dir
func getConfigDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		logger.Println("Failed to load user home dir: ", err)
		dir = "."
	}
	return filepath.Join(dir, KryptonConfigDir)
}

// return current profile's address setting for module/server
func (c *Config) GetAddress(module, server string) string {
	return settings.CurrentProfile.Server.getAddress(module, server)
}

// return current profile's token type setting
func (c *Config) GetTokenType() string {
	return settings.CurrentProfile.Auth.getTokenType()
}

// return current profile's management server setting
func (c *Config) GetManagementServer() string {
	return settings.CurrentProfile.Auth.getManagementServer()
}

// return cached access token if any
func (c *Config) GetAccessToken() string {
	return settings.EnrollTokenCache.GetAccessToken()
}

// return cached app token if any
func (c *Config) GetAppToken() string {
	return settings.AppTokenCache.GetAppToken()
}

// iot root cert path
func (c *Config) GetIotRootCertPath() string {
	return filepath.Join(getConfigDir(), IotRootCertFile)
}

// device token
func (c *Config) GetDeviceToken() string {
	return settings.DeviceTokenCache.GetDeviceToken()
}

// device id
func (c *Config) GetDeviceId() string {
	return settings.DeviceTokenCache.GetDeviceId()
}

// tenant id
func (c *Config) GetTenantId() string {
	return settings.DeviceTokenCache.GetTenantId()
}

// bulk enroll token
func (c *Config) GetBulkEnrollToken() string {
	return GetBulkEnrollCache().GetToken()
}

// bulk enroll token
func (c *Config) GetBulkEnrollTenantId() string {
	return GetBulkEnrollCache().GetTenantId()
}

// shortcut to set log level to verbose
func SetVerboseLogging() {
	logger.SetLevel(logging.Verbose)
}

// shortcut to set doc request
func SetDocType(docType string) {
	logger.SetDocType(docType)
}

func GetEncryptedTestData(name string) *EncryptedTestData {
	if testData, ok := settings.EncryptedTestData[name]; ok {
		return &testData
	}
	return nil
}
