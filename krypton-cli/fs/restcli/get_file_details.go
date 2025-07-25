package restcli

import (
	"cli/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FileDetailsClient struct {
	client *FilesClient
}

// {"file_id":25,"tenant_id":"1c815fef-0aea-484d-aadb-6b0debfbd8b5","device_id":"2c81a1f3-2b8d-4bd0-864a-3de4b3c2e92c","name":"fs_cli_upload_404991305","checksum":"L6vF+CaK/YjCxHInqgvT6A==","size":25,"status":"uploaded","created_at":"2023-05-04T03:41:14.720021Z","updated_at":"2023-05-04T03:41:14.765436Z"}
type FileDetailsResponseData struct {
	FileId    int    `json:"file_id"`
	Url       string `json:"url"`
	TenantId  string `json:"tenant_id"`
	DeviceId  string `json:"device_id"`
	Name      string `json:"name"`
	Checksum  string `json:"checksum"`
	Size      int    `json:"size"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type FileDetailsResponse struct {
	Data FileDetailsResponseData `json:"file"`
}

func NewFileDetailsClient(server, token string) *FileDetailsClient {
	var client = FileDetailsClient{
		client: filesClient(server, token, RetryCount),
	}
	return &client
}

func (c *FileDetailsClient) Execute(id int) (*FileDetailsResponse, error) {
	var resp *http.Response
	var err error
	url := c.client.getFileDetailsUrl(id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create file details url request. url: %s, Error: %v\n",
			url, err)
		return nil, err
	}
	common.AddUserAgentHeader(req)
	req.Header.Add("Authorization", "Bearer "+c.client.Token)
	client := common.RetriableClient(c.client.RetryCount)
	if resp, err = client.Do(req); err != nil {
		log.Printf("Failed to execute filedetaisl request. FileID: %d, Error: %v\n",
			id, err)
		return nil, err
	}
	return getFileDetailsResponse(resp)
}

func getFileDetailsResponse(resp *http.Response) (*FileDetailsResponse, error) {
	fd := FileDetailsResponse{}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get_file_details failed. Expected %d. Got %d",
			http.StatusOK, resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read file details response. Error: %v", err)
		return nil, err
	}
	log.HttpResponse(resp, data)
	err = json.Unmarshal(data, &fd)
	if err != nil {
		log.Printf("Failed to unmarshal file details response. Error: %v\n", err)
		return nil, err
	}
	return &fd, nil
}
