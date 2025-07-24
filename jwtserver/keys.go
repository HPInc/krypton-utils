package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Keys struct {
	Keys []jsonWebKey `json:"keys"`
}

func keysHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(Keys{Keys: GetJWKS()})
	if err != nil {
		log.Printf("Error serving keys: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
