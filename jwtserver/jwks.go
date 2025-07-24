package main

import (
	"encoding/base64"
	"math/big"
)

type jsonWebKey struct {
	Exponent string `json:"e"`
	ID       string `json:"kid"`
	Modulus  string `json:"n"`
	Type     string `json:"kty"`
	Use      string `json:"use"`
}

/*
	N:   b64.URLEncoding.EncodeToString(tokenVerificationKey.N.Bytes()),
	E: b64.URLEncoding.EncodeToString(common.NewBufferFromInt(
		uint64(tokenVerificationKey.E)).Data),
*/

func GetJWKS() []jsonWebKey {
	return []jsonWebKey{
		jsonWebKey{
			Exponent: getExponent(),
			ID:       KID,
			Modulus:  getModulus(),
			Type:     "RSA",
			Use:      "sig",
		},
	}
}

func getModulus() string {
	return base64.RawURLEncoding.EncodeToString(signKey.PublicKey.N.Bytes())
}

func getExponent() string {
	return base64.RawURLEncoding.EncodeToString(
		big.NewInt(int64(signKey.PublicKey.E)).Bytes())
}
