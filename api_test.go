package bepaid

import (
	"bepaid/vo"
	"context"
	"io"
	"net/http"
	"testing"
)

type A = vo.AuthorizationRequest

var ch = make(chan io.ReadCloser)

type customRoundTripper struct{}

func (c customRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	ch <- request.Body
	return nil, nil
}

func TestApi_Authorization(t *testing.T) {

	api := NewApi(
		&http.Client{Transport: customRoundTripper{}},
		DefaultEndpoints,
		",",
		"",
		"",
	)

	tests := []struct {
		name string
		data A
		er   string
	}{
		{"test1", A{}, `{"request":{"amount":0,"currency":"","tracking_id":"","test":false,"credit_card":{"number":"","verification_value":"","holder":"","exp_month":"","exp_year":"","skip_three_d_secure_verification":false}}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// or make ch buffered
			go api.Authorization(context.TODO(), tc.data)

			body := <-ch
			defer body.Close()

			b, err := io.ReadAll(body)

			if err != nil {
				t.Fatalf("wrong value: ER: %v, AR: %v", nil, err)
			}

			if string(b) != tc.er {
				t.Fatalf("wrong value: ER: %v, AR: %v", tc.er, string(b))
			}
		})
	}

}
