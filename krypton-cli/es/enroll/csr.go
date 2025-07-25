package enroll

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	keySize = 2048
)

// enroll payload consists of the following
// csr : get csr bytes and base64 encode it
// management service: this is set to HP connect with no config now
func (c *EnrollClient) getEnrollPayload() ([]byte, error) {
	csrBytes, err := c.createCSR()
	if err != nil {
		return nil, err
	}
	csr := base64.StdEncoding.EncodeToString(csrBytes)
	enrollData := enrollRequest{
		CSR:               csr,
		ManagementService: c.ManagementServer,
		HardwareHash:      c.HardwareHash,
	}
	jsonbytes, err := json.Marshal(enrollData)
	if err != nil {
		fmt.Printf("Failed to marshal enrollment payload. Error: %v\n", err)
		return nil, err
	}
	return jsonbytes, nil
}

// create a csr file and return the csrbytes
// note there is no pem encoding. raw csr bytes returned
func (c *EnrollClient) createCSR() ([]byte, error) {
	var err error
	c.PK, err = rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	deviceCsrTpl := x509.CertificateRequest{
		Subject:            pkix.Name{},
		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,
		PublicKey:          c.PK.PublicKey,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &deviceCsrTpl, c.PK)
	if err != nil {
		return nil, err
	}
	return csrBytes, nil
}
