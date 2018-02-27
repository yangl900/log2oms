package logclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	locationGMT = time.FixedZone("GMT", 0)
)

// LogClient is the client for log analytics
type LogClient struct {
	workspaceID     string
	workspaceSecret string
	logType         string
	httpClient      *http.Client
	signingKey      []byte
	apiLogsURL      string
}

type log struct {
	Data string `json:"data"`
}

// NewLogClient creates a log client
func NewLogClient(workspaceID, workspaceSecret, logType string) LogClient {
	client := LogClient{
		workspaceID:     workspaceID,
		workspaceSecret: workspaceSecret,
		logType:         logType,
	}

	client.httpClient = &http.Client{Timeout: time.Second * 30}
	client.signingKey, _ = base64.StdEncoding.DecodeString(workspaceSecret)
	client.apiLogsURL = fmt.Sprintf("https://%s.ods.opinsights.azure.com/api/logs?api-version=2016-04-01", workspaceID)

	return client
}

// PostMessage logs a single message to log analytics service
func (c *LogClient) PostMessage(message string) error {
	return c.PostMessages([]string{message})
}

// PostMessages logs an array of messages to log analytics service
func (c *LogClient) PostMessages(messages []string) error {
	var logs []log
	for _, m := range messages {
		logs = append(logs, log{Data: m})
	}

	body, _ := json.Marshal(logs)
	req, _ := http.NewRequest(http.MethodPost, c.apiLogsURL, bytes.NewReader(body))

	date := time.Now().In(locationGMT).Format(time.RFC1123)
	stringToSign := "POST\n" + strconv.FormatInt(req.ContentLength, 10) + "\napplication/json\n" + "x-ms-date:" + date + "\n/api/logs"

	signature := computeHmac256(stringToSign, c.signingKey)

	req.Header.Set("Authorization", fmt.Sprintf("SharedKey %s:%s", c.workspaceID, signature))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Log-Type", c.logType)
	req.Header.Set("x-ms-date", date)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to post request: %v", err)
	}

	if response.StatusCode != 200 {
		defer response.Body.Close()
		buf, _ := ioutil.ReadAll(response.Body)

		return fmt.Errorf("Post log request failed with status: %d %s", response.StatusCode, string(buf))
	}

	fmt.Printf("[%s] Posted %d messages.\n", date, len(logs))

	return nil
}
