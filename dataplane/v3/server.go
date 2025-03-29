package v3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Server struct {
	Id      *string `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
	Address *string `json:"address,omitempty"`
	Port    *int    `json:"port,omitempty"`
}

func (c Client) AddServer(backend string, transactionId string, b Server) (*Server, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/backends/%s/servers?transaction_id=%s",
		c.BaseUrl,
		backend,
		transactionId,
	)

	reqTxt, err := json.Marshal(b)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return c.executeApiReturnsServer(apiUrl, "POST", bytes.NewReader(reqTxt))
}

func (c Client) GetServer(name string, backend string, transactionId string) (*Server, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/backends/%s/servers/%s?transaction_id=%s",
		c.BaseUrl,
		backend,
		name,
		transactionId,
	)

	return c.executeApiReturnsServer(apiUrl, "GET", nil)
}

func (c Client) ListServer(backend string, transactionId string) ([]Server, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/backends/%s/servers?transaction_id=%s",
		c.BaseUrl,
		backend,
		transactionId,
	)

	resTxt, err := c.callApi(apiUrl, "GET", nil)
	if err != nil {
		return nil, err
	}

	if len(string(resTxt)) == 0 {
		return nil, nil
	}

	var resResult []Server
	if err := json.Unmarshal(resTxt, &resResult); err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return resResult, nil
}

func (c Client) ReplaceServer(backend string, transactionId string, b Server) (*Server, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/backends/%s/servers?transaction_id=%s",
		c.BaseUrl,
		backend,
		transactionId,
	)

	reqTxt, err := json.Marshal(b)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return c.executeApiReturnsServer(apiUrl, "PUT", bytes.NewReader(reqTxt))
}

func (c Client) DeleteServer(name string, backend string, transactionId string) error {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/backends/%s/servers/%s?transaction_id=%s",
		c.BaseUrl,
		backend,
		name,
		transactionId,
	)

	_, err := c.callApi(apiUrl, "DELETE", nil)

	return err
}

func (c Client) executeApiReturnsServer(apiUrl string, method string, body io.Reader) (*Server, error) {
	resTxt, err := c.callApi(apiUrl, method, body)
	if err != nil {
		return nil, err
	}

	if len(string(resTxt)) == 0 {
		return nil, nil
	}

	var resResult Server
	if err := json.Unmarshal(resTxt, &resResult); err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return &resResult, nil
}
