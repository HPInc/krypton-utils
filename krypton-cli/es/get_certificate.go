package es

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"cli/cmd"
	"cli/common"
	"cli/es/enroll"
)

const (
	CMD_GET_CERTIFICATE = "get_certificate"
)

type CmdGetCertificate struct {
	cmd.CmdBase
	EnrollBase
}

type CertDetails struct {
	DeviceId    string `json:"device_id"`
	Certificate string `json:"certificate"`
	Thumbprint  string `json:"thumbprint"`
	PrivateKey  string `json:"private_key"`
}

func init() {
	commands[CMD_GET_CERTIFICATE] = NewCmdGetCertificate()
}

func NewCmdGetCertificate() *CmdGetCertificate {
	c := CmdGetCertificate{
		cmd.CmdBase{Name: CMD_GET_CERTIFICATE},
		EnrollBase{},
	}
	fs := c.BaseInitFlags()
	(&c.EnrollBase).initFlags(fs)
	(&enrollFlags).initFlags(fs)
	return &c
}

func (c *CmdGetCertificate) Parse(args []string) (cmd.Command, error) {
	c.BaseParse(args)
	(&c.EnrollBase).initClient(c.RetryCount, c.ApiBasePath)
	c.Client.HardwareHash = *enrollFlags.hardwareHash
	c.Client.BulkEnrollToken = *enrollFlags.bulkEnrollToken
	c.RunFunc = c.getCertificate
	return c, nil
}

func (c *CmdGetCertificate) getCertificate() {
	client := c.Client
	cr, err := client.GetDeviceCertificate()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	details, err := getCertDetails(cr)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	details.PrivateKey = common.PemFromPrivateKey(client.PK)
	fmt.Println(common.GetJsonString(details))
}

func getCertDetails(cr *enroll.CertificateResponse) (*CertDetails, error) {
	cert, err := getx509Cert(&cr.Certificate)
	if err != nil {
		return nil, err
	}
	if err = verifyCertificate(cert); err != nil {
		return nil, err
	}
	if !verifyDeviceIdInCertificateCN(cert, cr.DeviceId) {
		return nil, errors.New("device id verification in cn failed")
	}
	return &CertDetails{
		DeviceId:    cr.DeviceId,
		Certificate: cr.Certificate,
		Thumbprint:  getCertificateThumbprint(cert),
	}, nil
}

func getx509Cert(encoded *string) (*x509.Certificate, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(*encoded)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(decodedBytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// Return a SHA256 checksum of the raw certificate as its thumbprint.
func getCertificateThumbprint(cert *x509.Certificate) string {
	thumbprint := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(thumbprint[:])
}

// VerifyCertificate - perform some verification checks on the certificate.
func verifyCertificate(cert *x509.Certificate) error {
	// Verify the certificate is currently valid and hasn't yet expired.
	if cert.NotBefore.After(time.Now()) {
		log.Error("cert nbf: ", cert.NotBefore, "current time: ", time.Now())
		return errors.New("certificate is not yet valid")
	}
	if cert.NotAfter.Before(time.Now()) {
		log.Error("cert naf: ", cert.NotAfter, "current time: ", time.Now())
		return errors.New("certificate has expired")
	}

	// Check the signature algorithm and public key algorithm.
	if cert.SignatureAlgorithm != x509.SHA256WithRSA {
		log.Error("cert sig alg: ", cert.SignatureAlgorithm)
		return errors.New("invalid signature algorithm")
	}
	if cert.PublicKeyAlgorithm != x509.RSA {
		log.Error("cert pubkey alg: ", cert.PublicKeyAlgorithm)
		return errors.New("invalid public key algorithm")
	}

	// Check if the key usage for the certificate is acceptable.
	if cert.KeyUsage != x509.KeyUsageDigitalSignature {
		log.Error("cert key usage: ", cert.KeyUsage)
		return errors.New("invalid key usage")
	}
	for _, usage := range cert.ExtKeyUsage {
		switch usage {
		case x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth:
			continue
		default:
			return errors.New("invalid ext key usage")
		}
	}
	return nil
}

// VerifyDeviceIdInCertificateCN - check if the device ID in the
// certificate's common name field matches the specified device ID.
func verifyDeviceIdInCertificateCN(cert *x509.Certificate, deviceID string) bool {
	return cert.Subject.CommonName == deviceID
}

func (c *CmdGetCertificate) GetInput() interface{} {
	return nil
}

func (c *CmdGetCertificate) ExecuteWithArgs(i interface{}) error {
	return nil
}
