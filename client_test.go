package bepaid

import (
	"bepaid/mocks"
	"bepaid/responses"
	"bepaid/vo"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"strings"
	"testing"
)

type responseBody struct {
	io.Reader
}

func (responseBody) Close() error {
	return nil
}

type TR = responses.TransactionResponse
type R = responses.Response

func TestClient_AuthorizationsWithResponse(t *testing.T) {

	ctrl := gomock.NewController(t)

	mockApi := mocks.NewMockApiInterface(ctrl)

	tests := []struct {
		name string
		data string
		er   R
	}{
		{"test1", `{ "response": { "message": "Unauthorized" } }`, R{Message: "Unauthorized", Errors: nil}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{Body: responseBody{strings.NewReader(tc.data)}}

			mockApi.EXPECT().Authorization(gomock.Any(), gomock.Any()).AnyTimes().Return(resp, nil)

			client := NewClient(mockApi)

			_, err := client.Authorizations(context.TODO(), *vo.NewAuthorizationRequest())

			if !cmp.Equal(err, &tc.er) {
				t.Fatalf("\ner: %T,\nar: %T,\ndiff:%s", &tc.er, err, cmp.Diff(err, &tc.er))
			}
		})
	}
}

func TestClient_AuthorizationsWithTransaction(t *testing.T) {

	ctrl := gomock.NewController(t)

	mockApi := mocks.NewMockApiInterface(ctrl)

	tests := []struct {
		name string
		data string
		er   TR
	}{
		{"", `{ "transaction": { "uid": "id123", "status":"successful", "amount": 123 } }`, TR{Uid: "id123", Status: "successful", Amount: 123}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{Body: responseBody{strings.NewReader(tc.data)}}

			mockApi.EXPECT().Authorization(gomock.Any(), gomock.Any()).AnyTimes().Return(resp, nil)

			client := NewClient(mockApi)

			ar, err := client.Authorizations(context.TODO(), *vo.NewAuthorizationRequest())
			if err != nil {
				t.Fatal(err)
			}

			if ar.Uid != tc.er.Uid || ar.Status != tc.er.Status || ar.Amount != tc.er.Amount {
				t.Fatal("not equal")
			}
		})
	}
}
