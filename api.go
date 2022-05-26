package bepaid

import (
	"bepaid/vo"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
)

//go:generate mockgen -source=api.go -destination=mocks/mock_api.go
type ApiInterface interface {
	Authorization(ctx context.Context, request vo.AuthorizationRequest) (*http.Response, error)
	Capture(ctx context.Context, capture vo.CaptureRequest) (*http.Response, error)
}

const (
	authorizations = "authorizations"
	captures       = "captures"
)

type Endpoint map[string]string

var DefaultEndpoints = map[string]string{
	authorizations: authorizations,
	captures:       captures,
}

type Api struct {
	client    *http.Client
	EndPoints Endpoint
	baseUrl   string
	auth      string
}

func NewApi(client *http.Client, endPoints Endpoint, baseUrl, username, password string) *Api {
	return &Api{
		client:    client,
		EndPoints: endPoints,
		baseUrl:   baseUrl,
		auth:      base64.StdEncoding.EncodeToString([]byte(username + ":" + password)),
	}
}

func (a *Api) Authorization(ctx context.Context, request vo.AuthorizationRequest) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodPost, authorizations, &request)
}

func (a *Api) Capture(ctx context.Context, request vo.CaptureRequest) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodPost, captures, &request)
}

func (a *Api) sendRequest(ctx context.Context, method, path string, request any) (*http.Response, error) {

	reader, err := marshalRequest(request)
	if err != nil {
		return nil, err
	}
	//body, err := json.Marshal(request)
	//if err != nil {
	//	return nil, err
	//}

	//r, err := http.NewRequestWithContext(ctx, method, a.baseUrl+path, bytes.NewReader(body))
	r, err := http.NewRequestWithContext(ctx, method, a.baseUrl+path, reader)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Basic "+a.auth)

	if method == http.MethodPost {
		//r.Header.Set("Content-Type", "application/json; charset=UTF-8")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept", "application/json")
	}

	return a.client.Do(r)
}

func marshalRequest(request any) (io.Reader, error) {
	b, err := json.Marshal(struct {
		Request any `json:"request"`
	}{request})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}