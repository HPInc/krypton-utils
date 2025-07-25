// wraps retry logic in overrides to http.Client
package common

import (
	"cli/logging"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	RetryAfter               = "Retry-After"
	DefaultRetryAfterSeconds = time.Second * 3
	DefaultRetryCount        = 3
)

type Do func(req *http.Request) (*http.Response, error)

type Client struct {
	HttpClient       *http.Client
	Do               Do
	MaxRetry         uint
	IgnoreStatusList []int
}

// no friction wrapper
func NewClient() *Client {
	return &Client{
		HttpClient: http.DefaultClient,
		Do:         http.DefaultClient.Do,
	}
}

// retry enabled client
func RetriableClient(retries uint) *Client {
	c := NewClient()
	c.Do = c.RetriableDo
	c.MaxRetry = retries
	return c
}

func (c *Client) Get(url string) (*http.Response, error) {
	return c.doRequest(url, "GET", nil)
}

// post with retry
func (c *Client) Post(url string, body io.Reader) (*http.Response, error) {
	return c.doRequest(url, "POST", body)
}

// handle requests
func (c *Client) doRequest(url, method string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	AddUserAgentHeader(req)
	logging.GetLogger().HttpRequest(req)
	return c.Do(req)
}

func (c *Client) RetriableDo(req *http.Request) (*http.Response, error) {
	logger := logging.GetLogger()
	logger.HttpRequest(req)
	var i uint
	for {
		resp, err := c.HttpClient.Do(req)
		if c.RetryWithBackoff(logger, i, resp, err) {
			logger.Debugf("%d/%d\n", i, c.MaxRetry)
			i++
		} else {
			return resp, err
		}
	}
}

func (c *Client) RetryWithBackoff(logger *logging.Log,
	i uint, resp *http.Response, err error) bool {
	doRetry := false
	retrySeconds := time.Second * time.Duration(i)
	if err != nil {
		logger.Errorf("error = %v, retriable.\n", err)
		doRetry = true
	} else if resp.StatusCode == 0 || resp.StatusCode > 500 {
		logger.Debugf("status = %d, retriable.\n", resp.StatusCode)
		doRetry = true
	} else if resp.StatusCode == 429 {
		retrySeconds = time.Duration(getRetryAfterHeaderValue(resp))
		logger.Debugf(
			"status = %d, retry_after = %d seconds, retriable.\n",
			resp.StatusCode, retrySeconds)
		doRetry = true
	} else if c.canIgnoreStatus(resp.StatusCode) {
		logger.Debugf("status = %d, retriable: configured as ignore.\n", resp.StatusCode)
		doRetry = true
	}
	if doRetry && i < c.MaxRetry {
		time.Sleep(retrySeconds)
		return true
	}
	return false
}

// check if client is configured to ignore a specific return status
// outside the common retriable ones
func (c *Client) canIgnoreStatus(status int) bool {
	for _, ignoreStatus := range c.IgnoreStatusList {
		if status == ignoreStatus {
			return true
		}
	}
	return false
}

func getRetryAfterHeaderValue(resp *http.Response) time.Duration {
	retryAfter := DefaultRetryAfterSeconds
	var err error
	var i int
	if val, ok := resp.Header[RetryAfter]; ok {
		if i, err = strconv.Atoi(val[0]); err != nil {
			logging.GetLogger().Errorf(
				"Error: invalid retry-after value: %v\n", val)
		}
		retryAfter = time.Second * time.Duration(i)
	}
	return retryAfter
}
