package contracts

import (
	"context"
	"github.com/TalismanFR/bepaid/service/vo"
	"net/http"
)

//go:generate mockgen -source=Api.go -destination=../../testdata/ApiMock.go -package=testdata
type Api interface {
	Payment(ctx context.Context, payment vo.PaymentRequest) (*http.Response, error)
	Authorization(ctx context.Context, authorization vo.AuthorizationRequest) (*http.Response, error)
	Capture(ctx context.Context, capture vo.CaptureRequest) (*http.Response, error)
	Void(ctx context.Context, void vo.VoidRequest) (*http.Response, error)
	Refund(ctx context.Context, refund vo.RefundRequest) (*http.Response, error)

	StatusByUid(ctx context.Context, uid string) (*http.Response, error)
	StatusByTrackingId(ctx context.Context, trackingId string) (*http.Response, error)
}
