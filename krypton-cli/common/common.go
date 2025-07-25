package common

import (
	"bytes"
	"cli/logging"
	"crypto/md5" // #nosec G501
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base32"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultContent = "stock default content"
	retryCount     = 3
)

type FileInfo struct {
	Name     string
	Checksum string
	Size     int64
}

func UploadFile(fileName string, data []byte, presignedURL string) error {
	log := logging.GetLogger()

	var f *os.File
	var err error
	if data == nil {
		f, err = OpenFile(fileName)
		if err != nil {
			return fmt.Errorf("Error opening file: %s", err)
		}
	}

	var reqBody io.ReadCloser
	if f != nil {
		reqBody = io.NopCloser(f)
	} else {
		reqBody = io.NopCloser(bytes.NewBuffer(data))
	}
	req, err := http.NewRequest(http.MethodPut, presignedURL, reqBody)
	if err != nil {
		return fmt.Errorf("Error making new request: %s", err)
	}

	var info *FileInfo
	if data != nil {
		info, err = GetDataInfo(data)
	} else {
		info, err = GetFileInfo(fileName)
	}
	if err != nil {
		return fmt.Errorf("Error getting file info: %s", err)
	}

	req.ContentLength = info.Size
	AddUserAgentHeader(req)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-MD5", info.Checksum)

	client := RetriableClient(DefaultRetryCount)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error executing request: %s", err)
	}
	defer resp.Body.Close()
	log.HttpResponse(resp, nil)
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Upload status: %d. Error: %s",
			resp.StatusCode,
			string(data))
	}
	return nil
}

// create temp file with random contents with a max length
// used for tests
func TempFileWithRandomContents(prefix string, maxLength int) (string, error) {
	file, err := os.CreateTemp("", prefix)
	if file != nil {
		defer file.Close()
		if _, err = file.Write(GetRandomBytes(maxLength)); err != nil {
			return "", err
		}
		return file.Name(), nil
	}
	return "", err
}

// get random bytes of english lowercase letters
// up to a max length
// note that errors are silently suppressed for default behavior
// not suitable for production.
func GetRandomBytes(maxLength int) []byte {
	n64Length := int64(maxLength)
	length, err := rand.Int(rand.Reader, big.NewInt(n64Length))
	if err != nil {
		length = big.NewInt(n64Length)
	}
	strLength := length.Uint64() + 1
	bytes := make([]byte, strLength)
	if _, err = rand.Read(bytes); err != nil {
		return []byte(defaultContent)
	}
	return []byte(strings.ToLower(base32.StdEncoding.EncodeToString(bytes)[:strLength]))
}

func GetFileInfo(filename string) (*FileInfo, error) {
	fs, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	checksum, err := getMd5Sum(filename)
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Name:     fs.Name(),
		Size:     fs.Size(),
		Checksum: checksum,
	}, nil
}

// return fileinfo struct from bytes
func GetDataInfo(data []byte) (*FileInfo, error) {
	md5bytes := md5.Sum(data) //#nosec G401
	return &FileInfo{
		Size:     int64(len(data)),
		Checksum: GetBase64(md5bytes[:]),
	}, nil
}

func getMd5Sum(file string) (string, error) {
	f, err := OpenFile(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// this is for s3 files. see krypton wiki for details.
	h := md5.New() //#nosec G401
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return GetBase64(h.Sum(nil)), nil
}

func GetBase64(data []byte) string {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)
	return string(dst)
}

// get private key as string
func PemFromPrivateKey(key *rsa.PrivateKey) string {
	bytes := x509.MarshalPKCS1PrivateKey(key)
	pemBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: bytes,
		},
	)
	return string(pemBytes)
}

// get rsa.PrivateKey from string
func PrivateKeyFromPem(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("pem decode failed")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func RetryWait(count uint, fn func() bool) bool {
	for i := uint(1); i <= count; i++ {
		if fn() {
			return true
		}
		logging.GetLogger().Debugf("Waiting %d seconds..", i)
		time.Sleep(time.Duration(i) * time.Second)
	}
	return false
}
