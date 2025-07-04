package v3

import (
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Credential string
	BaseUrl    string
}

type NormalResponse struct {
	Code    int     `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (c Client) constructAuthorizationHeader() http.Header {
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Basic %s", c.Credential))

	return h
}

func (c Client) callApi(apiUrl string, method string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, apiUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header = c.constructAuthorizationHeader()
	req.Header.Add("Content-Type", "application/json")

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resTxt, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("undefined response")
	}

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return resTxt, &UnauthorizedError{Message: string(resTxt)}
	case http.StatusBadRequest:
		return resTxt, &BadRequestError{Message: string(resTxt)}
	case http.StatusNotFound:
		return resTxt, &NotFoundError{Message: string(resTxt)}
	case http.StatusConflict:
		return resTxt, &ConflictError{Message: string(resTxt)}
	default:
		if res.StatusCode/100 != 2 { // 2xx status codes are successful
			return resTxt, &UnknownError{
				Message:    string(resTxt),
				StatusCode: res.StatusCode,
			}
		}
	}

	return resTxt, nil
}
