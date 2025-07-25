package common

import (
	"bytes"
	"cli/logging"
	"encoding/json"
)

func GetJsonString(r interface{}) string {
	buffer := &bytes.Buffer{}
	e := json.NewEncoder(buffer)
	e.SetEscapeHTML(false)
	err := e.Encode(r)
	if err != nil {
		logging.GetLogger().Fatal("could not encode json", err)
		return ""
	}
	return buffer.String()
}
