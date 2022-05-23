package bepaid

import (
	"bepaid/vo"
	"context"
	"encoding/base64"
	"io"
	"net/http"
)

//go:generate mockgen -source=api.go -destination=mocks/mock_api.go
type ApiInterface interface {
	Authorization(ctx context.Context, request vo.AuthorizationRequest) (*http.Response, error)
	Capture(ctx context.Context, capture vo.CaptureRequest) (*http.Response, error)
}

type Api struct {
	client  *http.Client
	baseUrl string
	auth    string
}

func NewApi(client *http.Client, baseUrl, username, password string) *Api {
	return &Api{client: client, baseUrl: baseUrl, auth: base64.StdEncoding.EncodeToString([]byte(username + ":" + password))}
}

func (a *Api) send(ctx context.Context, path, method string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequestWithContext(ctx, method, a.baseUrl+path, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Basic "+a.auth)
	return a.client.Do(r)
}

func (a *Api) Authorization(ctx context.Context, request vo.AuthorizationRequest) (*http.Response, error) {
	return a.send(ctx, "POST", "authorizations", request)
}

func (a *Api) Capture(ctx context.Context, request vo.CaptureRequest) (*http.Response, error) {
	return a.send(ctx, "POST", "captures", request)
}
