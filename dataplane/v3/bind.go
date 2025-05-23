package v3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Bind struct {
	Id      *string `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
	Address *string `json:"address,omitempty"`
	Port    *int    `json:"port,omitempty"`
	V4V6    *bool   `json:"v4v6,omitempty"`
	V6Only  *bool   `json:"v6only,omitempty"`
}

func (c Client) AddBind(frontend string, b Bind, opts ...option) (*Bind, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/frontends/%s/binds",
		c.BaseUrl,
		frontend,
	)

	op := &HAProxyOpts{}
	for _, o := range opts {
		o(op)
	}
	if op.TransactionID != "" {
		apiUrl += fmt.Sprintf("?transaction_id=%s", op.TransactionID)
	} else {
		return nil, &BadRequestError{Message: "transaction_id is required"}
	}

	reqTxt, err := json.Marshal(b)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return c.executeApiReturnsBind(apiUrl, "POST", bytes.NewReader(reqTxt))
}

func (c Client) GetBind(name string, frontend string, opts ...option) (*Bind, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/frontends/%s/binds/%s",
		c.BaseUrl,
		frontend,
		name,
	)

	op := &HAProxyOpts{}
	for _, o := range opts {
		o(op)
	}
	if op.TransactionID != "" {
		apiUrl += fmt.Sprintf("?transaction_id=%s", op.TransactionID)
	}

	return c.executeApiReturnsBind(apiUrl, "GET", nil)
}

func (c Client) ListBind(frontend string, opts ...option) ([]Bind, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/frontends/%s/binds",
		c.BaseUrl,
		frontend,
	)

	op := &HAProxyOpts{}
	for _, o := range opts {
		o(op)
	}

	if op.TransactionID != "" {
		apiUrl += fmt.Sprintf("?transaction_id=%s", op.TransactionID)
	}

	resTxt, err := c.callApi(apiUrl, "GET", nil)
	if err != nil {
		return nil, err
	}

	if len(string(resTxt)) == 0 {
		return nil, nil
	}

	var resResult []Bind
	if err := json.Unmarshal(resTxt, &resResult); err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return resResult, nil
}

func (c Client) ReplaceBind(frontend string, b Bind, opts ...option) (*Bind, error) {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/frontends/%s/binds",
		c.BaseUrl,
		frontend,
	)

	op := &HAProxyOpts{}
	for _, o := range opts {
		o(op)
	}

	if op.TransactionID != "" {
		apiUrl += fmt.Sprintf("?transaction_id=%s", op.TransactionID)
	} else {
		return nil, &BadRequestError{Message: "transaction_id is required"}
	}

	reqTxt, err := json.Marshal(b)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return c.executeApiReturnsBind(apiUrl, "PUT", bytes.NewReader(reqTxt))
}

func (c Client) DeleteBind(name string, frontend string, opts ...option) error {
	apiUrl := fmt.Sprintf(
		"%s/v3/services/haproxy/configuration/frontends/%s/binds/%s",
		c.BaseUrl,
		frontend,
		name,
	)

	op := &HAProxyOpts{}
	for _, o := range opts {
		o(op)
	}

	if op.TransactionID != "" {
		apiUrl += fmt.Sprintf("?transaction_id=%s", op.TransactionID)
	} else {
		return &BadRequestError{Message: "transaction_id is required"}
	}

	_, err := c.callApi(apiUrl, "DELETE", nil)

	return err
}

func (c Client) executeApiReturnsBind(apiUrl string, method string, body io.Reader) (*Bind, error) {
	resTxt, err := c.callApi(apiUrl, method, body)
	if err != nil {
		return nil, err
	}

	if len(string(resTxt)) == 0 {
		return nil, nil
	}

	var resResult Bind
	if err := json.Unmarshal(resTxt, &resResult); err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return &resResult, nil
}
