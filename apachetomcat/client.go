package apachetomcat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ApacheTomcatStatus struct {
	Tomcat struct {
		JVM struct {
			Memory struct {
				Free  string `json:"free"`
				Total string `json:"total"`
				Max   string `json:"max"`
			} `json:"memory"`
			MemoryPool []struct {
				Name           string `json:"name"`
				Type           string `json:"type"`
				UsageInit      string `json:"usageInit"`
				UsageCommitted string `json:"usageCommitted"`
				UsageMax       string `json:"usageMax"`
				UsageUsed      string `json:"usageUsed"`
			} `json:"memorypool"`
		} `json:"jvm"`
		Connector []struct {
			Name       string `json:"name"`
			ThreadInfo struct {
				MaxThreads         string `json:"maxThreads"`
				CurrentThreadCount string `json:"currentThreadCount"`
				CurrentThreadsBusy string `json:"currentThreadsBusy"`
			} `json:"threadInfo"`
			RequestInfo struct {
				MaxTime        string `json:"maxTime"`
				ProcessingTime string `json:"processingTime"`
				RequestCount   string `json:"requestCount"`
				ErrorCount     string `json:"errorCount"`
				BytesReceived  string `json:"bytesReceived"`
				BytesSent      string `json:"bytesSent"`
			} `json:"requestInfo"`
		} `json:"connector"`
	} `json:"tomcat"`
}

type apacheTomcatClient struct {
	client http.Client
	config *Config
}

func newDefaultApacheTomcatClient(config *Config) *apacheTomcatClient {
	return &apacheTomcatClient{
		client: http.Client{
			// TODO: let user specify this timeout?
			Timeout: time.Second * 60,
		},
		config: config,
	}
}

func (c *apacheTomcatClient) getTomcatStatus() (ApacheTomcatStatus, error) {
	managerStatusPath := "/manager/status/all?JSON=true"
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", c.config.Endpoint, managerStatusPath), http.NoBody)
	req.SetBasicAuth(c.config.Username, c.config.Password)

	var apacheTomcatStatus ApacheTomcatStatus

	// TODO: should we be using fmt.Sprintf here?
	resp, err := c.client.Do(req)
	if err != nil {
		return apacheTomcatStatus, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return apacheTomcatStatus, err
	}

	err = json.Unmarshal(respBody, &apacheTomcatStatus)
	if err != nil {
		return apacheTomcatStatus, err
	}

	return apacheTomcatStatus, nil
}

/*

Metrics:
	MEMORY:
		attributes: NONE
		tomcat.jvm.memory.free
		tomcat.jvm.memory.total
		tomcat.jvm.memory.max

	MEMORY POOL:
		attributes:
		- tomcat.jvm.memory_pool.name
		- tomcat.jvm.memory_pool.type
		tomcat.jvm.memory_pool.usage_init
		tomcat.jvm.memory_pool.usage_committed
		tomcat.jvm.memory_pool.usage_max
		tomcat.jvm.memory_pool.usage_used

	CONNECTOR:
		attributes:
		- tomcat.connector.name
		tomcat.connector.thread_info.max_threads
		tomcat.connector.thread_info.current_thread_count
		tomcat.connector.thread_info.current_thread_busy
		tomcat.connector.request_info.max_time
		tomcat.connector.request_info.processing_time
		tomcat.connector.request_info.request_count
		tomcat.connector.request_info.error_count
		tomcat.connector.request_info.bytes_received
		tomcat.connector.request_info.bytes_sent
*/
