package v3

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	TRANSACTION_STATUS_FAILED      = "failed"
	TRANSACTION_STATUS_OUTDATED    = "outdated"
	TRANSACTION_STATUS_IN_PROGRESS = "in_progress"
	TRANSACTION_STATUS_SUCCESS     = "success"
)

type Transaction struct {
	Id     *string `json:"id,omitempty"`
	Status *string `json:"status,omitempty"`
}

func (c Client) GetVersion() (*int, error) {
	apiUrl := fmt.Sprintf("%s/services/haproxy/configuration/version", c.BaseUrl)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, &InternalError{Message: err.Error()}
	}

	req.Header = c.constructAuthorizationHeader()
	req.Header.Add("Content-Type", "application/json")
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	versionString := strings.TrimRight(string(result), "\n")
	version, err := strconv.Atoi(versionString)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	return &version, nil
}

func (c Client) CreateTransaction(version int) (*Transaction, error) {
	apiUrl := fmt.Sprintf("%s/services/haproxy/transactions?version=%d", c.BaseUrl, version)

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return nil, &InternalError{Message: err.Error()}
	}

	req.Header = c.constructAuthorizationHeader()
	req.Header.Add("Content-Type", "application/json")
	return c.executeApiReturnsTransaction(req)
}

func (c Client) GetTransaction(id string) (*Transaction, error) {
	apiUrl := fmt.Sprintf("%s/services/haproxy/transactions/%s", c.BaseUrl, id)

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return nil, &InternalError{Message: err.Error()}
	}

	req.Header = c.constructAuthorizationHeader()
	req.Header.Add("Content-Type", "application/json")
	return c.executeApiReturnsTransaction(req)
}

func (c Client) CommitTransaction(id string) (*Transaction, error) {
	apiUrl := fmt.Sprintf("%s/services/haproxy/transactions/%s", c.BaseUrl, id)

	req, err := http.NewRequest("PUT", apiUrl, nil)
	if err != nil {
		return nil, &InternalError{Message: err.Error()}
	}

	req.Header = c.constructAuthorizationHeader()
	req.Header.Add("Content-Type", "application/json")
	return c.executeApiReturnsTransaction(req)
}

func (c Client) CloseTransaction(id string) (*string, error) {
	apiUrl := fmt.Sprintf("%s/services/haproxy/transactions/%s", c.BaseUrl, id)
	res, err := c.callApi(apiUrl, "DELETE", nil)
	if err != nil {
		return nil, err
	}

	responseText := string(res)
	return &responseText, err
}

func (c Client) executeApiReturnsTransaction(r *http.Request) (*Transaction, error) {
	client := new(http.Client)
	res, err := client.Do(r)
	if err != nil {
		return nil, &InternalError{Message: err.Error()}
	}

	resTxt, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &InvalidResponseError{Message: err.Error()}
	}

	if len(string(resTxt)) == 0 {
		return nil, nil
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, &UnauthorizedError{Message: "unauthorized"}
	} else if res.StatusCode == http.StatusNotFound {
		return nil, &NotFoundError{Message: "not found"}
	} else if res.StatusCode/100 != 2 {
		return nil, &UnknownError{Message: "unknown error"}
	}

	var resResult Transaction
	if err := json.Unmarshal(resTxt, &resResult); err != nil {
		return nil, &InvalidResponseError{
			Message: err.Error(),
		}
	}

	return &resResult, nil
}
