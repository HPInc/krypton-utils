package main

import (
	"log"
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
		//Server
		"JWT_PORT": {v: &c.Server.Port},
		//TOKEN
		"JWT_TOKEN_AUDIENCE":      {v: &c.Token.Audience},
		"JWT_TOKEN_ISSUER":        {v: &c.Token.Issuer},
		"JWT_TOKEN_VALID_MINUTES": {v: &c.Token.ValidMinutes},
	}
	for k, v := range m {
		e := os.Getenv(k)
		if e != "" {
			log.Printf("Overriding env variable: %s = %s\n",
				k, getLoggableValue(v.secret, e))
			replaceConfigValue(os.Getenv(k), &v)
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
			log.Printf("Error: Bad integer value in env: %s", envValue)
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
