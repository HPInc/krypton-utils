package restcli

import (
	"cli/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type DownloadUrlClient struct {
	client *FilesClient
}

type DownloadUrlRequest struct {
	FileName string `json:"name"`
	Checksum string `json:"checksum"`
	Size     int64  `json:"size"`
}

type DownloadUrlResponse struct {
	FileName string `json:"file_name"`
	FileId   int    `json:"file_id"`
	Url      string `json:"url"`
}

func NewDownloadUrlClient(server, token string) *DownloadUrlClient {
	var client = DownloadUrlClient{
		client: filesClient(server, token, RetryCount),
	}
	return &client
}

func (c *DownloadUrlClient) Execute(id int) (*DownloadUrlResponse, error) {
	downloadUrl := c.client.getDownloadUrl(id)
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		log.Printf("Failed to create downloadurl request. Error: %v\n", err)
		return nil, err
	}
	client := common.RetriableClient(RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send enrollment request. Retries: %d, Error: %v\n",
			RetryCount, err)
		return nil, err
	}
	return getDownloadUrlResponse(id, resp)
}

func getDownloadUrlResponse(id int, resp *http.Response) (*DownloadUrlResponse, error) {
	ur := DownloadUrlResponse{FileId: id}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetFile failed. Expected %d. Got %d\n",
			http.StatusOK, resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read downloadurl response. Error: %v", err)
		return nil, err
	}
	log.HttpResponse(resp, data)
	err = json.Unmarshal(data, &ur)
	if err != nil {
		log.Printf("Failed to unmarshal downloadurl response. Error: %v", err)
		return nil, err
	}
	return &ur, nil
}

func (u *DownloadUrlResponse) GetDownloadUrl() (string, error) {
	if u.Url != "" {
		return u.Url, nil
	}
	return "", errors.New("download url response did not have a valid url")
}
