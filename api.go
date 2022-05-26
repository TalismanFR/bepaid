package bepaid

import (
	"bepaid/vo"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
)

//go:generate mockgen -source=api.go -destination=mocks/mock_api.go -package=mocks
type ApiInterface interface {
	Authorization(ctx context.Context, request vo.AuthorizationRequest) (*http.Response, error)
	Capture(ctx context.Context, capture vo.CaptureRequest) (*http.Response, error)
}

const (
	authorizations = "authorizations"
	captures       = "captures"
)

type Api struct {
	client  *http.Client
	baseUrl string
	auth    string
}

func NewApi(client *http.Client, baseUrl, username, password string) *Api {
	return &Api{
		client:  client,
		baseUrl: baseUrl,
		auth:    base64.StdEncoding.EncodeToString([]byte(username + ":" + password)),
	}
}

func (a *Api) Authorization(ctx context.Context, request vo.AuthorizationRequest) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodPost, authorizations, &request)
}

func (a *Api) Capture(ctx context.Context, request vo.CaptureRequest) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodPost, captures, &request)
}

func (a *Api) sendRequest(ctx context.Context, method, path string, request any) (*http.Response, error) {

	r, w := io.Pipe()
	go func() {
		json.NewEncoder(w).Encode(struct {
			Request any `json:"request"`
		}{request})
		w.Close()
	}()

	req, err := http.NewRequestWithContext(ctx, method, a.baseUrl+path, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+a.auth)
	req.Header.Set("Accept", "application/json")

	if method == http.MethodPost {
		//r.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("Content-Type", "application/json")
	}

	return a.client.Do(req)
}
