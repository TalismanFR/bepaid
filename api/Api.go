package api

import (
	"bepaid-sdk/service/vo"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const (
	payments         = "/transactions/payments"
	authorizations   = "/transactions/authorizations"
	captures         = "/transactions/captures"
	voids            = "/transactions/voids"
	refunds          = "/transactions/refunds"
	statusUid        = "/transactions/"
	statusTrackingId = "/v2/transactions/tracking_id/"
)

type Api struct {
	client  *http.Client
	baseUrl string
	auth    string
}

func NewApi(client *http.Client, baseUrl, username, password string) *Api {
	return &Api{
		client:  client,
		baseUrl: strings.TrimSuffix(baseUrl, "/"),
		auth:    "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))}
}

func (a *Api) StatusByUid(ctx context.Context, uid string) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodGet, statusUid+uid, nil)
}

func (a *Api) StatusByTrackingId(ctx context.Context, trackingId string) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodGet, statusTrackingId+trackingId, nil)
}

func (a *Api) Payment(ctx context.Context, payment vo.PaymentRequest) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodPost, payments, &payment)
}

func (a *Api) Authorization(ctx context.Context, authorization vo.AuthorizationRequest) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodPost, authorizations, &authorization)
}

func (a *Api) Capture(ctx context.Context, capture vo.CaptureRequest) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodPost, captures, &capture)
}

func (a *Api) Void(ctx context.Context, void vo.VoidRequest) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodPost, voids, &void)
}

func (a *Api) Refund(ctx context.Context, refund vo.RefundRequest) (*http.Response, error) {
	return a.sendRequest(ctx, http.MethodPost, refunds, &refund)
}

func (a *Api) sendRequest(ctx context.Context, method, path string, request interface{}) (*http.Response, error) {

	//if request == nil {
	//
	//}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, method, a.baseUrl+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", a.auth)
	r.Header.Set("Accept", "application/json")

	if method == http.MethodPost {
		//r.Header.Set("Content-Type", "application/json; charset=UTF-8")
		r.Header.Set("Content-Type", "application/json")
	}

	return a.client.Do(r)
}

// if request doesnt have "request" field
func marshalRequest(request interface{}) (io.Reader, error) {
	b, err := json.Marshal(struct {
		Request interface{} `json:"request"`
	}{request})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (a Api) GetUrl() string {
	return a.baseUrl
}
