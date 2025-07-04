package v3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const (
	BACKEND_BALANCE_ALGORITHM_FIRST      = "first"
	BACKEND_BALANCE_ALGORITHM_HASH       = "hash"
	BACKEND_BALANCE_ALGORITHM_RANDOM     = "random"
	BACKEND_BALANCE_ALGORITHM_ROUNDROBIN = "roundrobin"
)

const (
	BACKEND_MODE_TCP = "tcp"
)

type BackendBalance struct {
	Algorithm string `json:"algorithm,omitempty"`
}

type Backend struct {
	Id      *int            `json:"id,omitempty"`
	Balance *BackendBalance `json:"balance,omitempty"`
	Name    *string         `json:"name,omitempty"`
	Mode    string          `json:"mode,omitempty"`
}

func (c Client) AddBackend(backend Backend, transactionId string) (*Backend, error) {
	apiUrl := fmt.Sprintf("%s/v3/services/haproxy/configuration/backends?transaction_id=%s", c.BaseUrl, transactionId)

	body, err := json.Marshal(backend)
	if err != nil {
		return nil, &InvalidResponseError{
			Message: err.Error(),
		}
	}

	return c.executeApiReturnsBackend(apiUrl, "POST", bytes.NewReader(body))
}

func (c Client) GetBackend(name string, transactionId string) (*Backend, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/backends/%s?transaction_id=%s",
		c.BaseUrl,
		name,
		transactionId,
	)

	return c.executeApiReturnsBackend(apiUrl, "GET", nil)
}

func (c Client) ListBackend(transactionId string) ([]Backend, error) {
	apiUrl := fmt.Sprintf("%s/v3/services/haproxy/configuration/backends?transaction_id=%s", c.BaseUrl, transactionId)

	resTxt, err := c.callApi(apiUrl, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("invalid response payload")
	}

	if len(string(resTxt)) == 0 {
		return nil, nil
	}

	var resResult []Backend
	if err := json.Unmarshal(resTxt, &resResult); err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return resResult, nil
}

func (c Client) ReplaceBackend(name string, backend Backend, transactionId string) (*Backend, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/backends/%s?transaction_id=%s",
		c.BaseUrl,
		name,
		transactionId,
	)

	body, err := json.Marshal(backend)
	if err != nil {
		return nil, err
	}

	return c.executeApiReturnsBackend(apiUrl, "PUT", bytes.NewReader(body))
}

func (c Client) DeleteBackend(name string, transactionId string) error {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/backends/%s?transaction_id=%s",
		c.BaseUrl,
		name,
		transactionId,
	)
	_, err := c.callApi(apiUrl, "DELETE", nil)

	return err
}

func (c Client) executeApiReturnsBackend(apiUrl string, method string, body io.Reader) (*Backend, error) {
	resTxt, err := c.callApi(apiUrl, method, body)
	if err != nil {
		return nil, err
	}

	if len(string(resTxt)) == 0 {
		return nil, nil
	}

	var resResult Backend
	if err := json.Unmarshal(resTxt, &resResult); err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return &resResult, nil
}
