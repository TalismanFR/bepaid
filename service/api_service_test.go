package service

import (
	"bepaid-sdk/service/vo"
	"bepaid-sdk/testdata"
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
)

const (
	json_req_1 = `{
	   "transaction":{
	      "capture":{
	         "message":"The operation was successfully processed.",
	         "ref_id":"8889999",
	         "gateway_id":152,
	         "status":"successful"
	      },
	      "uid":"2-310b0da80b",
	      "status":"successful",
	      "message":"Successfully processed",
	      "amount":50,
	      "parent_uid":"1-310b0da80b",
	      "receipt_url": "",
	      "currency":"USD",
	      "type":"capture",
	      "test":false
	   }
	}`
)

func TestApiService_Capture(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ioutil.NopCloser(bytes.NewReader([]byte(json_req_1)))
	capture := testdata.NewMockApi(ctrl)

	capture.EXPECT().Capture(context.Background(), *vo.NewCaptureRequest(50, "1-310b0da80b")).Return(http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)

	captureTest := NewApiService(capture)
	response, err := captureTest.Capture(context.Background(), *vo.NewCaptureRequest(50, "1-310b0da80b"))

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "successful", response.Transaction.Status)
	assert.Equal(t, 50, response.Transaction.Amount)
	assert.Equal(t, "1-310b0da80b", response.Transaction.ParentUid)
}

func TestApiService_Capture2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ioutil.NopCloser(bytes.NewReader([]byte(json_req_1)))
	capture := testdata.NewMockApi(ctrl)

	capture.EXPECT().Capture(context.Background(), *vo.NewCaptureRequest(50, "1-310b0da80b")).Return(http.Response{
		StatusCode: 200,
		Body:       r,
	}, errors.New("error message"))

	captureTest := NewApiService(capture)
	_, err := captureTest.Capture(context.Background(), *vo.NewCaptureRequest(50, "1-310b0da80b"))

	assert.NotNil(t, err)
	assert.Equal(t, "error message", err.Error())
}

func TestApiService_Capture3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ioutil.NopCloser(bytes.NewReader([]byte(json_req_1)))
	capture := testdata.NewMockApi(ctrl)

	capture.EXPECT().Capture(context.Background(), *vo.NewCaptureRequest(50, "1-310b0da80b")).Return(http.Response{
		StatusCode: 100,
		Body:       r,
	}, nil)

	captureTest := NewApiService(capture)
	_, err := captureTest.Capture(context.Background(), *vo.NewCaptureRequest(50, "1-310b0da80b"))

	assert.NotNil(t, err)
	assert.Equal(t, "Error, status code: 100", err.Error())
}
