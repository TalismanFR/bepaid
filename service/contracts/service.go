package contracts

import (
	"context"
	"github.com/TalismanFR/bepaid/service/vo"
)

type ApiService interface {
	Authorizations(ctx context.Context, authorizationRequest vo.AuthorizationRequest) (vo.TransactionResponse, error)
	Capture(ctx context.Context, captureRequest vo.CaptureRequest) (vo.TransactionResponse, error)
}
