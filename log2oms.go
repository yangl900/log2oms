package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/yangl900/log2oms/logclient"
)

const (
	envLogFile         = "LOG2OMS_LOG_FILE"
	envLogType         = "LOG2OMS_LOG_TYPE"
	envWorkspaceID     = "LOG2OMS_WORKSPACE_ID"
	envWorkspaceSecret = "LOG2OMS_WORKSPACE_SECRET"
	envMetadataPrefix  = "LOG2OMS_METADATA_"
)

var (
	batchSizeInLines = 100000
	requestSizeLimit = 1024 * 1024 * 8
)

func logLines(client *logclient.LogClient, lines []string) {
	err := client.PostMessages(lines, time.Now().UTC())
	if err != nil {
		fmt.Println(err)
	}
}

func metadata() map[string]string {
	metadata := make(map[string]string)
	metadata["Hostname"], _ = os.Hostname()

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")

		if strings.HasPrefix(pair[0], envMetadataPrefix) {
			key := strings.TrimPrefix(pair[0], envMetadataPrefix)
			metadata[key] = pair[1]
		}
	}

	return metadata
}

func main() {
	workspaceID, workspaceSecret := os.Getenv(envWorkspaceID), os.Getenv(envWorkspaceSecret)
	if workspaceID == "" || workspaceSecret == "" {
		fmt.Printf("Workspace Id and secret not defined in environment variable '%s' and '%s'\n", envWorkspaceID, envWorkspaceSecret)
		return
	}

	logfile := os.Getenv(envLogFile)
	if logfile == "" {
		if len(os.Args) < 2 {
			fmt.Printf("Neither '%s' environment variable nor command line parameter specified.\n", envLogFile)
			return
		}

		logfile = os.Args[1]
	}

	logType := os.Getenv(envLogType)
	if logType == "" {
		logType = "container_logs"
	}

	metadata := metadata()
	for m := range metadata {
		fmt.Printf("[LOG2OMS][%s] %s = %s\n", time.Now().UTC().Format(time.RFC3339), m, metadata[m])
	}

	fmt.Printf("[LOG2OMS][%s] Start tail logs from: %s\n", time.Now().UTC().Format(time.RFC3339), logfile)

	client := logclient.NewLogClient(workspaceID, workspaceSecret, logType, metadata)

	t, err := tail.TailFile(logfile, tail.Config{ReOpen: true, Follow: true})
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := []string{}
	byteCount := 0
	for {
		select {
		case line := <-t.Lines:
			if line.Err != nil {
				fmt.Println(line.Err)
			} else {
				fmt.Printf("[%s] %s\n", line.Time.UTC().Format(time.RFC3339), line.Text)
			}

			lines = append(lines, line.Text)
			byteCount += len(line.Text)

			if len(lines) >= batchSizeInLines || byteCount >= requestSizeLimit {
				logLines(&client, lines)
				lines = []string{}
				byteCount = 0
			}
		case <-time.After(time.Second * 5):
			if len(lines) > 0 {
				logLines(&client, lines)
				lines = []string{}
				byteCount = 0
			}
		}
	}
}
