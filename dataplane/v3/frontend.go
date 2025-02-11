package v3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const (
	FRONTEND_MODE_TCP = "tcp"
)

type Frontend struct {
	DefaultBackend *string `json:"default_backend,omitempty"`
	Description    *string `json:"description,omitempty"`
	Disabled       *bool   `json:"disabled,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty"`
	Id             *int    `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	Mode           *string `json:"mode"`
}

func (c Client) AddFrontend(f Frontend, transactionId string) (*Frontend, error) {
	apiUrl := fmt.Sprintf("%s/services/haproxy/configuration/frontends?transaction_id=%s", c.BaseUrl, transactionId)

	body, err := json.Marshal(f)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return c.executeApiReturnsFrontend(apiUrl, "POST", bytes.NewReader(body))
}

func (c Client) GetFrontend(name string, transactionId string) (*Frontend, error) {
	apiUrl := fmt.Sprintf(
		"%s/services/haproxy/configuration/frontends/%s?transaction_id=%s",
		c.BaseUrl,
		name,
		transactionId,
	)

	return c.executeApiReturnsFrontend(apiUrl, "GET", nil)
}

func (c Client) ListFrontend(transactionId string) ([]Frontend, error) {
	apiUrl := fmt.Sprintf("%s/services/haproxy/configuration/frontends?transaction_id=%s", c.BaseUrl, transactionId)

	resTxt, err := c.callApi(apiUrl, "GET", nil)
	if err != nil {
		return nil, err
	}

	if len(string(resTxt)) == 0 {
		return nil, nil
	}

	var resResult []Frontend
	if err := json.Unmarshal(resTxt, &resResult); err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return resResult, nil
}

func (c Client) ReplaceFrontend(name string, f Frontend, transactionId string) (*Frontend, error) {
	apiUrl := fmt.Sprintf(
		"%s/services/haproxy/configuration/frontends/%s?transaction_id=%s",
		c.BaseUrl,
		name,
		transactionId,
	)

	body, err := json.Marshal(f)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return c.executeApiReturnsFrontend(apiUrl, "PUT", bytes.NewReader(body))
}

func (c Client) DeleteFrontend(name string, transactionId string) error {
	apiUrl := fmt.Sprintf(
		"%s/services/haproxy/configuration/frontends/%s?transaction_id=%s",
		c.BaseUrl,
		name,
		transactionId,
	)

	_, err := c.callApi(apiUrl, "DELETE", nil)

	return err
}

func (c Client) executeApiReturnsFrontend(apiUrl string, method string, body io.Reader) (*Frontend, error) {
	resTxt, err := c.callApi(apiUrl, method, body)
	if err != nil {
		return nil, err
	}

	if len(string(resTxt)) == 0 {
		return nil, nil
	}

	var resResult Frontend
	if err := json.Unmarshal(resTxt, &resResult); err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return &resResult, nil
}
