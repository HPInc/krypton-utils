package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	envPort          = "WEBHOOK_PORT"
	envSQSHost       = "SQS_HOST"
	envSQSPort       = "SQS_PORT"
	envSQSQueueName  = "SQS_QUEUE_NAME"
	envSQSRetryAfter = "SQS_RETRY_AFTER"
)

var (
	port          = 9001
	sqsHost       = "localhost"
	sqsPort       = 9324
	sqsUrl        = ""
	sqsQueueName  = "queue"
	sqsRetryAfter = 3
)

func main() {
	getConfig()
	sqsUrl = fmt.Sprintf("http://%s:%d/queue/%s", sqsHost, sqsPort,
		sqsQueueName)

	http.HandleFunc("/", webhookHandler)

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// forward incoming webhook notifications to sqs
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s3EventData, _ := io.ReadAll(r.Body)
	err := relayMessageToSQS(s3EventData)
	if err != nil {
		w.Header().Add("Retry-After", fmt.Sprintf("%d", sqsRetryAfter))
		w.WriteHeader(http.StatusTooManyRequests)
	}
}

func relayMessageToSQS(msg []byte) error {
	data := fmt.Sprintf("Action=SendMessage&MessageBody=%s",
		strings.Replace(string(msg), " ", "%20", -1))
	req, err := http.NewRequest(http.MethodPost, sqsUrl,
		strings.NewReader(data))
	if err != nil {
		log.Printf("Could not create relay request: %v\n", err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("queue forward error: %d", resp.StatusCode)
	}
	return nil
}

// read config vals
func getConfig() {
	port = getIntEnv(envPort)
	sqsPort = getIntEnv(envSQSPort)
	sqsHost = getStringEnv(envSQSHost, sqsHost)
	sqsQueueName = getStringEnv(envSQSQueueName, sqsQueueName)
	sqsRetryAfter = getIntEnv(envSQSRetryAfter)
}

// get string env val or default
func getStringEnv(name, defaultVal string) string {
	strVal := os.Getenv(name)
	if strVal == "" {
		strVal = defaultVal
	}
	return strVal
}

// get int env val
func getIntEnv(name string) int {
	var intval int
	var err error
	envValue := os.Getenv(name)
	if envValue != "" {
		if intval, err = strconv.Atoi(envValue); err != nil {
			log.Fatalf("Invalid value for env: %s. %v", envPort, err)
		}
	}
	return intval
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
