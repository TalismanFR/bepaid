// Code generated by MockGen. DO NOT EDIT.
// Source: api.go

// Package mock_bepaid is a generated GoMock package.
package mocks

import (
	vo "bepaid/vo"
	context "context"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockApiInterface is a mock of ApiInterface interface.
type MockApiInterface struct {
	ctrl     *gomock.Controller
	recorder *MockApiInterfaceMockRecorder
}

// MockApiInterfaceMockRecorder is the mock recorder for MockApiInterface.
type MockApiInterfaceMockRecorder struct {
	mock *MockApiInterface
}

// NewMockApiInterface creates a new mock instance.
func NewMockApiInterface(ctrl *gomock.Controller) *MockApiInterface {
	mock := &MockApiInterface{ctrl: ctrl}
	mock.recorder = &MockApiInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApiInterface) EXPECT() *MockApiInterfaceMockRecorder {
	return m.recorder
}

// Authorization mocks base method.
func (m *MockApiInterface) Authorization(ctx context.Context, request vo.AuthorizationRequest) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorization", ctx, request)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorization indicates an expected call of Authorization.
func (mr *MockApiInterfaceMockRecorder) Authorization(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorization", reflect.TypeOf((*MockApiInterface)(nil).Authorization), ctx, request)
}

// Capture mocks base method.
func (m *MockApiInterface) Capture(ctx context.Context, capture vo.CaptureRequest) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capture", ctx, capture)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Capture indicates an expected call of Capture.
func (mr *MockApiInterfaceMockRecorder) Capture(ctx, capture interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capture", reflect.TypeOf((*MockApiInterface)(nil).Capture), ctx, capture)
}