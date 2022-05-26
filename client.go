package bepaid

import (
	"bepaid/responses"
	"bepaid/vo"
	"context"
	"encoding/json"
)

//go:generate mockgen -source=client.go -destination=mocks/mock_client.go -package=mocks
type ClientInterface interface {
	Authorizations(ctx context.Context, authorizationRequest vo.AuthorizationRequest) (*responses.TransactionResponse, error)
	Capture(ctx context.Context, captureRequest vo.CaptureRequest) (*responses.TransactionResponse, error)
}

type Client struct {
	api ApiInterface
}

func NewClient(api ApiInterface) *Client {
	return &Client{api: api}
}

func (c *Client) Authorizations(ctx context.Context, authorizationRequest vo.AuthorizationRequest) (*responses.TransactionResponse, error) {
	response, err := c.api.Authorization(ctx, authorizationRequest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	v := struct {
		Transaction responses.TransactionResponse `json:"transaction"`
		Response    responses.Response            `json:"response"`
	}{}

	//b ,_ := io.ReadAll(response.Body)
	//fmt.Println(string(b))
	err = json.NewDecoder(response.Body).Decode(&v)
	if err != nil {
		return nil, err
	}

	if v.Response.Message != "" {
		return nil, &v.Response
	}

	return &v.Transaction, nil
}

func (c *Client) Capture(ctx context.Context, captureRequest vo.CaptureRequest) (*responses.TransactionResponse, error) {
	response, err := c.api.Capture(ctx, captureRequest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	v := struct {
		Transaction responses.TransactionResponse `json:"transaction"`
		Response    responses.Response            `json:"response"`
	}{}

	//b ,_ := io.ReadAll(response.Body)
	//fmt.Println(string(b))
	err = json.NewDecoder(response.Body).Decode(&v)
	if err != nil {
		return nil, err
	}

	if v.Response.Message != "" {
		return nil, &v.Response
	}

	return &v.Transaction, nil
}
