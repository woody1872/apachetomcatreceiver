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
		} `json:"jvm"`
	} `json:"tomcat"`
}

type apacheTomcatClient struct {
	client http.Client
	config *Config
}

func newDefaultApacheTomcatClient(config *Config) *apacheTomcatClient {
	return &apacheTomcatClient{
		client: http.Client{
			// TODO: read timeout from collectors config.yaml?
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
