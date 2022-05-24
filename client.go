package bepaid

import (
	"bepaid/response"
	"bepaid/vo"
	"context"
	"io"
	"os"
)

//go:generate mockgen -source=client.go -destination=mocks/mock_client.go
type ClientInterface interface {
	Authorizations(ctx context.Context, authorizationRequest vo.AuthorizationRequest) (*response.TransactionResponse, error)
	Capture(ctx context.Context, captureRequest vo.CaptureRequest) (*response.TransactionResponse, error)
}

type Client struct {
	api ApiInterface
}

func NewClient(api ApiInterface) *Client {
	return &Client{api: api}
}

func (c *Client) Authorizations(ctx context.Context, authorizationRequest vo.AuthorizationRequest) (*response.TransactionResponse, error) {
	response, err := c.api.Authorization(ctx, authorizationRequest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (c *Client) Capture(ctx context.Context, captureRequest vo.CaptureRequest) (*response.TransactionResponse, error) {
	response, err := c.api.Capture(ctx, captureRequest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		return nil, err
	}

	return nil, err
}
