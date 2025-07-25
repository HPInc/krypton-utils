package config

import (
	"os"
	"strconv"
)

type value struct {
	secret bool
	v      interface{}
}

// loadEnvironmentVariableOverrides - check values specified for supported
// environment variables. These can be used to override configuration settings
// specified in the config file.
func (c *Config) OverrideFromEnvironment() {
	m := map[string]value{
		"CLI_PROFILE_NAME": {v: &c.ProfileName},
	}
	for k, v := range m {
		e := os.Getenv(k)
		if e != "" {
			logger.Debugf("Overriding env variable: %s with %s\n",
				k, getLoggableValue(v.secret, e))
			val := v
			replaceConfigValue(os.Getenv(k), &val)
		}
	}
}

// envValue will be non empty as this function is private to file
func replaceConfigValue(envValue string, t *value) {
	switch t.v.(type) {
	case *string:
		*t.v.(*string) = envValue
	case *int:
		i, err := strconv.Atoi(envValue)
		if err != nil {
			logger.Error("Bad integer value in env")
		} else {
			*t.v.(*int) = i
		}
	}
}

func getLoggableValue(secret bool, value string) string {
	if secret {
		return "***"
	}
	return value
}
