package restcli

import (
	"cli/logging"
	"fmt"
)

type FilesClient struct {
	Server     string
	Token      string
	UploadUrl  string
	RetryCount uint
}

const (
	API_PATH          = "api/v1/files"
	INTERNAL_API_PATH = "api/internal/v1/files"
)

var log = logging.GetLogger()

func filesClient(server, token string, retryCount uint) *FilesClient {
	return &FilesClient{
		Server:     server,
		Token:      token,
		UploadUrl:  fmt.Sprintf("%s/%s", server, API_PATH),
		RetryCount: retryCount,
	}
}

func (c *FilesClient) getDownloadUrl(id int) string {
	return fmt.Sprintf("%s/%s/%d/signed_url?method=get",
		c.Server, INTERNAL_API_PATH, id)
}

func (c *FilesClient) getFileDetailsUrl(id int) string {
	return fmt.Sprintf("%s/%s/%d",
		c.Server, API_PATH, id)
}
