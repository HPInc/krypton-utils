package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	API_VERSION = "/api/v1/"
	privKeyPath = "privateKey.pem"
)

var (
	signKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
	config  Config = Config{}
	KID            = uuid.New().String()
)

type AADClaims struct {
	*jwt.StandardClaims
	TenantId string `json:"tid"`
}

type DeviceTokenClaims struct {
	*jwt.StandardClaims
	TenantId string `json:"tid"`
	Type     string `json:"typ"`
}

func main() {
	getConfig()
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", config.Server.Port)

	fmt.Println(API_VERSION + "token")
	http.HandleFunc(API_VERSION+"token", getTokenHandler)
	http.HandleFunc(API_VERSION+"keys", keysHandler)
	http.HandleFunc(API_VERSION+"device_token", getDeviceTokenHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getConfig() {
	config.Server.Port = 9090
	config.Token.Audience = "https://graph.microsoft.com"
	config.Token.Issuer = "https://sts.windows.net"
	config.Token.ValidMinutes = 5

	config.OverrideFromEnvironment()
}

func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	handleError(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func createToken(tenantId string) (string, error) {
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Header["kid"] = KID

	// set our claims
	t.Claims = &AADClaims{
		&jwt.StandardClaims{
			Audience: config.Token.Audience,
			Issuer:   config.Token.Issuer,
			Subject:  tenantId,
			ExpiresAt: time.Now().Add(
				time.Minute * time.Duration(config.Token.ValidMinutes)).Unix(),
		},
		tenantId,
	}

	// Creat token string
	return t.SignedString(signKey)
}

func getTokenHandler(w http.ResponseWriter, req *http.Request) {
	// check if tenant_id is specified as a url param
	tenantId := req.URL.Query().Get("tenant_id")
	if tenantId == "" {
		tenantId = uuid.New().String()
	}
	tokenString, err := createToken(tenantId)
	if err != nil {
		log.Printf("Token Signing error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error signing token!")
		return
	}
	w.Header().Set("Content-Type", "application/jwt")
	fmt.Fprintln(w, tokenString)
}

// mimic a device token
func getDeviceTokenHandler(w http.ResponseWriter, req *http.Request) {
	// check if tenant_id is specified as a url param
	tenantId := req.URL.Query().Get("tenant_id")
	if tenantId == "" {
		tenantId = uuid.New().String()
	}
	deviceId := req.URL.Query().Get("device_id")
	if deviceId == "" {
		deviceId = uuid.New().String()
	}
	audience := req.URL.Query().Get("audience")
	if audience == "" {
		audience = config.Token.Audience
	}
	issuer := req.URL.Query().Get("issuer")
	if issuer == "" {
		issuer = "HP Device Token Service"
	}
	tokenType := req.URL.Query().Get("type")
	if tokenType == "" {
		tokenType = "device"
	}
	tokenString, err := createDeviceToken(tenantId, deviceId, audience, issuer, tokenType)
	if err != nil {
		log.Printf("Token Signing error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error signing token!")
		return
	}
	w.Header().Set("Content-Type", "application/jwt")
	fmt.Fprintln(w, tokenString)
}

func createDeviceToken(tenantId, deviceId, audience, issuer, tokenType string) (string, error) {
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Header["kid"] = KID

	// set our claims
	t.Claims = &DeviceTokenClaims{
		&jwt.StandardClaims{
			Audience: audience,
			Issuer:   issuer,
			Subject:  deviceId,
			ExpiresAt: time.Now().Add(
				time.Minute * time.Duration(config.Token.ValidMinutes)).Unix(),
		},
		tenantId,
		tokenType,
	}

	// Creat token string
	return t.SignedString(signKey)
}
