package restcli

import (
	"bytes"
	"cli/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	RetryMultiplier = 2
	RetryCount      = 3
)

type UploadUrlClient struct {
	client *FilesClient
}

type UploadUrlParams struct {
	FileName string
	FileData []byte
}

type UploadUrlRequest struct {
	FileName string `json:"name"`
	Checksum string `json:"checksum"`
	Size     int64  `json:"size"`
}

type UploadUrlResponse struct {
	Data FileDetailsResponseData `json:"file"`
}

func NewUploadUrlClient(server, token string) *UploadUrlClient {
	var client = UploadUrlClient{
		client: filesClient(server, token, RetryCount),
	}
	return &client
}

func (c *UploadUrlClient) Execute(p UploadUrlParams) (*UploadUrlResponse, error) {
	jsonString, err := p.getUploadUrlPayload()
	if err != nil {
		log.Printf("Failed to initialize uploadurl payload. Error: %v", err)
		return nil, err
	}
	req, err := http.NewRequest(
		"POST", c.client.UploadUrl, bytes.NewBuffer(jsonString))
	if err != nil {
		log.Printf("Failed to create uploadurl request. Error: %v", err)
		return nil, err
	}
	common.AddUserAgentHeader(req)
	req.Header.Add("Authorization", "Bearer "+c.client.Token)
	req.Header.Add("Content-Type", "application/json")
	client := common.RetriableClient(RetryCount)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute uploadurl request, Error: %v\n", err)
		return nil, err
	}
	return getUploadUrlResponse(resp)
}

func getUploadUrlResponse(resp *http.Response) (*UploadUrlResponse, error) {
	ur := UploadUrlResponse{}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CreateFile failed. Expected %d. Got %d\n",
			http.StatusCreated, resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read uploadurl response. Error: %v\n", err)
		return nil, err
	}
	log.HttpResponse(resp, data)
	err = json.Unmarshal(data, &ur)
	if err != nil {
		log.Printf("Failed to unmarshal uploadurl response. Error: %v\n", err)
		return nil, err
	}
	return &ur, nil
}

func (p *UploadUrlParams) getUploadUrlPayload() ([]byte, error) {
	var fi *common.FileInfo
	var err error
	if p.FileData != nil {
		fi, err = common.GetDataInfo(p.FileData)
	} else {
		fi, err = common.GetFileInfo(p.FileName)
	}
	if err != nil {
		return nil, err
	}
	if fi.Name == "" {
		fi.Name = p.FileName
	}
	req := UploadUrlRequest{
		FileName: fi.Name,
		Checksum: fi.Checksum,
		Size:     fi.Size,
	}
	jsonbytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal uploadurl payload. Error: %v\n", err)
		return nil, err
	}
	return jsonbytes, nil
}

func (u *UploadUrlResponse) GetUploadUrl() (string, error) {
	if u.Data.FileId > 0 {
		return u.Data.Url, nil
	}
	return "", errors.New("upload url response did not have a valid url")
}
