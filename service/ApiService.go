package service

import (
	"bepaid-sdk/api/contracts"
	"bepaid-sdk/service/vo"
	"context"
)

type ApiService struct {
	api contracts.Api
}

func NewApiService(api contracts.Api) *ApiService {
	return &ApiService{api: api}
}

func (a ApiService) Authorizations(ctx context.Context, authorizationRequest vo.AuthorizationRequest) (vo.TransactionResponse, error) {
	//TODO implement me
	panic("implement me")

}

func (a ApiService) Capture(ctx context.Context, captureRequest vo.CaptureRequest) (vo.TransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a ApiService) CancelPayment(ctx context.Context, captureRequest vo.CaptureRequest) (vo.TransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a ApiService) ReturnPayment(ctx context.Context, captureRequest vo.CaptureRequest) (vo.TransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}
