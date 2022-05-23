package bepaid

import (
	"bepaid/vo"
	"context"
)

//go:generate mockgen -source=client.go -destination=mocks/mock_client.go
type ClientInterface interface {
	Authorizations(ctx context.Context, authorizationRequest vo.AuthorizationRequest) (*vo.TransactionResponse, error)
	Capture(ctx context.Context, captureRequest vo.CaptureRequest) (*vo.TransactionResponse, error)
}

type Client struct {
	api ApiInterface
}

func NewClient(api ApiInterface) *Client {
	return &Client{api: api}
}

func (a Client) Authorizations(ctx context.Context, authorizationRequest vo.AuthorizationRequest) (*vo.TransactionResponse, error) {
	//TODO implement me
	panic("implement me")

}

func (a Client) Capture(ctx context.Context, captureRequest vo.CaptureRequest) (*vo.TransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}
