package api

import (
	"bepaid-sdk/api/contracts"
	"bepaid-sdk/service/vo"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"io"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

type P = vo.PaymentRequest
type A = vo.AuthorizationRequest
type C = vo.CaptureRequest
type V = vo.VoidRequest
type R = vo.RefundRequest

var (

	// correctApi sends requests to real host and return host's response.
	//
	// correctApi initialised in TestMain method.
	//
	// Host address, shop_id, secret retrieved from .env file in 'bepaid' directory.
	correctApi *Api
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		panic("cannot load .env file")
	}

	correctApi = NewApi(http.DefaultClient, os.Getenv("API_HOST"), os.Getenv("SHOP_ID"), os.Getenv("SECRET"))

	m.Run()
}

func TestMarshalling(t *testing.T) {
	s := MarshalRequestTestSuite{}
	api := &Api{client: s.GetHttpClient()}
	s.Setup(api)
	s.Run(t)
	s.TearDown()
}

func TestSingleRequests(t *testing.T) {
	s := &SingleRequestSuite{}
	s.Setup(correctApi)
	suite.Run(t, s)
}

func TestSequentialRequests(t *testing.T) {
	s := &SequentialRequestsTestSuite{}
	s.Setup(correctApi)
	suite.Run(t, s)
}

//////////////////////////////////////
//	Marshalling request body suite	//
//////////////////////////////////////

type customRoundTripper struct {
	ch chan io.ReadCloser
}

func (c *customRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	c.ch <- request.Body
	return nil, nil
}

type MarshalRequestTestSuite struct {
	ch  chan io.ReadCloser
	api contracts.Api
}

// GetHttpClient return http.Client, that should be inserted in the implementation of contracts.Api
//
// Then use Setup to insert instance of that implementation.
func (s *MarshalRequestTestSuite) GetHttpClient() *http.Client {
	// tests, using mockApi, read Request.Body from this channel. So, no parallel testing.
	s.ch = make(chan io.ReadCloser, 1)

	// testingClient sends Request.Body to ch.
	testingClient := &http.Client{Transport: &customRoundTripper{s.ch}}

	// mockApi sends Request.Body to ch.
	// Use to check Request.Body marshalling
	return testingClient
}

func (s *MarshalRequestTestSuite) Setup(api contracts.Api) {
	s.api = api
}

func (s *MarshalRequestTestSuite) TearDown() {
	close(s.ch)
}

func (s *MarshalRequestTestSuite) Run(t *testing.T) {
	if s.api == nil {
		panic("api value is nil. Run Setup first.")
	}
	t.Run("Payment", func(t *testing.T) {
		s.Payment(t)
	})
	t.Run("Authorization", func(t *testing.T) {
		s.Authorization(t)
	})
	t.Run("Capture", func(t *testing.T) {
		s.Capture(t)
	})
	t.Run("Void", func(t *testing.T) {
		s.Void(t)
	})
	t.Run("Refund", func(t *testing.T) {
		s.Refund(t)
	})
}

func (s *MarshalRequestTestSuite) Payment(t *testing.T) {

	tests := []struct {
		name string
		req  P
		er   string
	}{
		{"defaultValue", P{}, `{"request":{"amount":0,"currency":"","description":"","tracking_id":"","test":false,"credit_card":{"number":"","verification_value":"","holder":"","exp_month":"","exp_year":"","skip_three_d_secure_verification":false}}}`},
		{"requestConstructor", *vo.NewPaymentRequest(vo.Amount(1), vo.RUB, "rub_1", "id1", true, *vo.NewCreditCard("5555", "123", "tim", "05", "2024")), `{"request":{"amount":1,"currency":"RUB","description":"rub_1","tracking_id":"id1","test":true,"credit_card":{"number":"5555","verification_value":"123","holder":"tim","exp_month":"05","exp_year":"2024","skip_three_d_secure_verification":false}}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				s.ch,
				tc.er,
				func() (*http.Response, error) {
					return s.api.Payment(context.TODO(), tc.req)
				})
		})
	}
}

func (s *MarshalRequestTestSuite) Authorization(t *testing.T) {

	tests := []struct {
		name string
		req  A
		er   string
	}{
		{"defaultValue", A{}, `{"request":{"amount":0,"currency":"","description":"","tracking_id":"","test":false,"credit_card":{"number":"","verification_value":"","holder":"","exp_month":"","exp_year":"","skip_three_d_secure_verification":false}}}`},
		{"requestConstructor", *vo.NewAuthorizationRequest(vo.Amount(1), vo.RUB, "rub_1", "id1", true, *vo.NewCreditCard("5555", "123", "tim", "05", "2024")), `{"request":{"amount":1,"currency":"RUB","description":"rub_1","tracking_id":"id1","test":true,"credit_card":{"number":"5555","verification_value":"123","holder":"tim","exp_month":"05","exp_year":"2024","skip_three_d_secure_verification":false}}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				s.ch,
				tc.er,
				func() (*http.Response, error) {
					return s.api.Authorization(context.TODO(), tc.req)
				})
		})
	}
}

func (s *MarshalRequestTestSuite) Capture(t *testing.T) {

	tests := []struct {
		name string
		req  vo.CaptureRequest
		er   string
	}{
		{"defaultValue", C{}, `{"request":{"parent_uid":"","amount":0}}`},
		{"requestConstructor", *vo.NewCaptureRequest("id123", vo.Amount(63)), `{"request":{"parent_uid":"id123","amount":63}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				s.ch,
				tc.er,
				func() (*http.Response, error) {
					return s.api.Capture(context.TODO(), tc.req)
				})
		})
	}
}

func (s *MarshalRequestTestSuite) Void(t *testing.T) {

	tests := []struct {
		name string
		req  V
		er   string
	}{
		{"defaultValue", V{}, `{"request":{"parent_uid":"","amount":0}}`},
		{"requestConstructor", *vo.NewVoidRequest("id123", vo.Amount(63)), `{"request":{"parent_uid":"id123","amount":63}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				s.ch,
				tc.er,
				func() (*http.Response, error) {
					return s.api.Void(context.TODO(), tc.req)
				})
		})
	}
}

func (s *MarshalRequestTestSuite) Refund(t *testing.T) {

	tests := []struct {
		name string
		req  R
		er   string
	}{
		{"defaultValue", R{}, `{"request":{"parent_uid":"","amount":0,"reason":""}}`},
		{"requestConstructor", *vo.NewRefundRequest("id123", vo.Amount(63), "reason"), `{"request":{"parent_uid":"id123","amount":63,"reason":"reason"}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				s.ch,
				tc.er,
				func() (*http.Response, error) {
					return s.api.Refund(context.TODO(), tc.req)
				})
		})
	}
}

func testMarshallRequest(t *testing.T, ch chan io.ReadCloser, er string, startRequest func() (*http.Response, error)) {
	// ignore response and error
	go startRequest()

	body := <-ch
	defer body.Close()

	b, err := io.ReadAll(body)

	if err != nil {
		fatalfWithExpectedActual(t, "ReadAll returned not nil value", nil, err)
	}

	if string(b) != er {
		fatalfWithExpectedActual(t, "Strings aren't equal", er, string(b))
	}
}

//////////////////////////////////////
// 		Single requests	 suite		//
//////////////////////////////////////

type SingleRequestSuite struct {
	suite.Suite
	api contracts.Api
}

func (s *SingleRequestSuite) Setup(api contracts.Api) {
	s.api = api
}

func (s *SingleRequestSuite) TestApi_Authorization() {
	tests := []struct {
		name   string
		a      A
		code   int
		status string
	}{
		{"test1",
			*vo.NewAuthorizationRequest(101, vo.RUB, "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
			http.StatusOK,
			"successful",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.T().Logf("UID: %s", apiAuthorization(s.T(), s.api, tc.a, tc.code, tc.status))
		})
	}
}

func (s *SingleRequestSuite) TestApi_Payment() {

	tests := []struct {
		name   string
		p      P
		code   int
		status string
	}{
		{"test1",
			*vo.NewPaymentRequest(vo.Amount(100), vo.RUB, "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")).WithDuplicateCheck(false),
			http.StatusOK,
			"successful",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.T().Logf("UID: %s", apiPayment(s.T(), s.api, tc.p, tc.code, tc.status))
		})
	}
}

func (s *SingleRequestSuite) TestApi_StatusByTrackingId() {
	r, _ := s.api.StatusByTrackingId(context.Background(), "mytrackingid")

	defer r.Body.Close()
	buf := bytes.Buffer{}
	b, _ := io.ReadAll(r.Body)
	json.Indent(&buf, b, "", "\t")
}

func (s *SingleRequestSuite) TestApi_StatusByUid() {
	r, _ := s.api.StatusByUid(context.Background(), "151534003-9d0e9c9aa1")

	defer r.Body.Close()
	buf := bytes.Buffer{}
	b, _ := io.ReadAll(r.Body)
	json.Indent(&buf, b, "", "\t")
}

//////////////////////////////////////
//		Sequential requests	suite	//
//////////////////////////////////////

type SequentialRequestsTestSuite struct {
	suite.Suite
	api contracts.Api
}

func (s *SequentialRequestsTestSuite) Setup(api contracts.Api) {
	s.api = api
}

func (s *SequentialRequestsTestSuite) TestApi_AuthorizationCapture() {

	type (
		Auth struct {
			a      A
			code   int
			status string
		}
		Capt struct {
			cWithUId func(string) C
			code     int
			status   string
		}
	)
	amount := vo.Amount(rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100)

	tests := []struct {
		name string
		auth Auth
		capt Capt
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, vo.RUB, "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			Capt{
				func(uid string) C { return *vo.NewCaptureRequest(uid, amount) },
				http.StatusOK,
				"successful",
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {

			uid := apiAuthorization(s.T(), correctApi, tc.auth.a, tc.auth.code, tc.auth.status)
			s.T().Logf("A.Uid: %s", uid)

			c := tc.capt.cWithUId(uid)

			uid = apiCapture(s.T(), correctApi, c, tc.capt.code, tc.capt.status)
			s.T().Logf("C.Uid: %s", uid)
		})
	}

}

func (s *SequentialRequestsTestSuite) TestApi_AuthorizationVoid() {
	type (
		Auth struct {
			a      A
			code   int
			status string
		}
		Void struct {
			vWithUid func(string) V
			code     int
			status   string
		}
	)
	amount := vo.Amount(rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100)

	tests := []struct {
		name string
		auth Auth
		void Void
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, vo.RUB, "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			Void{
				func(uid string) V { return *vo.NewVoidRequest(uid, amount) },
				http.StatusOK,
				"successful"},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {

			uid := apiAuthorization(s.T(), correctApi, tc.auth.a, tc.auth.code, tc.auth.status)
			s.T().Logf("A.Uid: %s", uid)

			v := tc.void.vWithUid(uid)

			uid = apiVoid(s.T(), correctApi, v, tc.void.code, tc.void.status)
			s.T().Logf("V.Uid: %s", uid)
		})
	}
}

func (s *SequentialRequestsTestSuite) TestApi_AuthorizationCaptureRefund() {
	type (
		Auth struct {
			a      A
			code   int
			status string
		}
		Capt struct {
			cWithUId func(string) C
			code     int
			status   string
		}
		Refund struct {
			rWithUid func(string) R
			code     int
			status   string
		}
	)
	amount := vo.Amount(rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100)

	tests := []struct {
		name   string
		auth   Auth
		capt   Capt
		refund Refund
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, vo.RUB, "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			Capt{
				func(uid string) C { return *vo.NewCaptureRequest(uid, amount) },
				http.StatusOK,
				"successful"},
			Refund{
				func(uid string) R { return *vo.NewRefundRequest(uid, amount, "give me my money back") },
				http.StatusOK,
				"successful",
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {

			uid := apiAuthorization(s.T(), correctApi, tc.auth.a, tc.auth.code, tc.auth.status)
			s.T().Logf("A.Uid: %s", uid)

			c := tc.capt.cWithUId(uid)

			uid = apiCapture(s.T(), correctApi, c, tc.capt.code, tc.capt.status)
			s.T().Logf("C.Uid: %s", uid)

			r := tc.refund.rWithUid(uid)

			uid = apiRefund(s.T(), correctApi, r, tc.refund.code, tc.refund.status)
			s.T().Logf("R.Uid: %s", uid)
		})
	}
}

func (s *SequentialRequestsTestSuite) TestApi_PaymentRefund() {
	type (
		Paym struct {
			p      P
			code   int
			status string
		}
		Refund struct {
			rWithUid func(string) R
			code     int
			status   string
		}
	)
	amount := vo.Amount(rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100)

	tests := []struct {
		name   string
		paym   Paym
		refund Refund
	}{
		{"positiveTest1",
			Paym{
				*vo.NewPaymentRequest(amount, vo.RUB, "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			Refund{
				func(uid string) R { return *vo.NewRefundRequest(uid, amount, "give me my money back") },
				http.StatusOK,
				"successful"},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {

			uid := apiPayment(s.T(), correctApi, tc.paym.p, tc.paym.code, tc.paym.status)
			s.T().Logf("A.Uid: %s", uid)

			r := tc.refund.rWithUid(uid)

			uid = apiRefund(s.T(), correctApi, r, tc.refund.code, tc.refund.status)
			s.T().Logf("R.Uid: %s", uid)
		})
	}
}

func (s *SequentialRequestsTestSuite) TestApi_AuthorizationStatusByUid() {
	type (
		Auth struct {
			a      A
			code   int
			status string
		}
		StatusUid struct {
			uid    string
			code   int
			status string
		}
	)
	amount := vo.Amount(rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100)

	tests := []struct {
		name    string
		auth    Auth
		statUid StatusUid
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, vo.RUB, "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			StatusUid{
				"",
				http.StatusOK,
				"successful",
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {

			uid := apiAuthorization(s.T(), correctApi, tc.auth.a, tc.auth.code, tc.auth.status)
			s.T().Logf("A.Uid: %s", uid)

			tc.statUid.uid = uid

			uid = apiStatusByUid(s.T(), correctApi, tc.statUid.uid, tc.statUid.code, tc.statUid.status)
			s.T().Logf("S.Uid: %s", uid)
		})
	}
}

func (s *SequentialRequestsTestSuite) TestApi_AuthorizationStatusByTrackingId() {
	type (
		Auth struct {
			a      A
			code   int
			status string
		}
		StatusTrId struct {
			trackingId string
			code       int
			status     string
		}
	)
	amount := vo.Amount(rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100)

	tests := []struct {
		name       string
		auth       Auth
		statusTrId StatusTrId
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, vo.RUB, "it's description", "mytrackingid5678", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			StatusTrId{
				"mytrackingid5678",
				http.StatusOK,
				"successful",
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {

			uid := apiAuthorization(s.T(), correctApi, tc.auth.a, tc.auth.code, tc.auth.status)
			s.T().Logf("A.Uid: %s", uid)

			uid = apiStatusByTrackingId(s.T(), correctApi, tc.statusTrId.trackingId, tc.statusTrId.code, tc.statusTrId.status)
			s.T().Logf("S.Uid: %s", uid)
		})
	}
}

//////////////////////////////////////
//		Helper functions			//
//////////////////////////////////////

func apiPayment(t *testing.T, api contracts.Api, request vo.PaymentRequest, codeExp int, statusExp string) string {
	resp, err := api.Payment(context.Background(), request)
	if err != nil {
		t.Fatal(sprintfExpAct("api.Authorization: return non nil error", nil, err))
	}
	if resp.StatusCode != codeExp {
		t.Fatal(sprintfExpAct("Response.StatusCode: unexpected status codeExp", codeExp, resp.StatusCode))
	}

	transaction, err := getTransaction(resp.Body)
	if err != nil {
		t.Fatal(sprintfExpAct("getTransaction: return non nil error", nil, err))
	}

	status, ok := transaction["status"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'status'", "key 'status' present", "no such key"))
	}

	statusS, ok := status.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'status'", "string", fmt.Sprintf("%T", status)))
	}

	if statusS != statusExp {
		t.Fatal(sprintfExpAct("transaction: status not equal", statusExp, statusS))
	}

	uid, ok := transaction["uid"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'uid'", "key 'uid' present", "no such key"))
	}

	uidS, ok := uid.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'uid'", "string", fmt.Sprintf("%T", uid)))
	}

	return uidS
}

func apiAuthorization(t *testing.T, api contracts.Api, request vo.AuthorizationRequest, codeExp int, statusExp string) string {
	resp, err := api.Authorization(context.Background(), request)
	if err != nil {
		t.Fatal(sprintfExpAct("api.Authorization: return non nil error", nil, err))
	}
	if resp.StatusCode != codeExp {
		t.Fatal(sprintfExpAct("Response.StatusCode: unexpected status codeExp", codeExp, resp.StatusCode))
	}

	transaction, err := getTransaction(resp.Body)
	if err != nil {
		t.Fatal(sprintfExpAct("getTransaction: return non nil error", nil, err))
	}

	status, ok := transaction["status"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'status'", "key 'status' present", "no such key"))
	}

	statusS, ok := status.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'status'", "string", fmt.Sprintf("%T", status)))
	}

	if statusS != statusExp {
		t.Fatal(sprintfExpAct("transaction: status not equal", statusExp, statusS))
	}

	uid, ok := transaction["uid"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'uid'", "key 'uid' present", "no such key"))
	}

	uidS, ok := uid.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'uid'", "string", fmt.Sprintf("%T", uid)))
	}

	return uidS
}

func apiCapture(t *testing.T, api contracts.Api, request vo.CaptureRequest, codeExp int, statusExp string) string {
	resp, err := api.Capture(context.Background(), request)
	if err != nil {
		t.Fatal(sprintfExpAct("api.Authorization: return non nil error", nil, err))
	}
	if resp.StatusCode != codeExp {
		t.Fatal(sprintfExpAct("Response.StatusCode: unexpected status codeExp", codeExp, resp.StatusCode))
	}

	transaction, err := getTransaction(resp.Body)
	if err != nil {
		t.Fatal(sprintfExpAct("getTransaction: return non nil error", nil, err))
	}

	status, ok := transaction["status"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'status'", "key 'status' present", "no such key"))
	}

	statusS, ok := status.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'status'", "string", fmt.Sprintf("%T", status)))
	}

	if statusS != statusExp {
		t.Fatal(sprintfExpAct("transaction: status not equal", statusExp, statusS))
	}

	uid, ok := transaction["uid"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'uid'", "key 'uid' present", "no such key"))
	}

	uidS, ok := uid.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'uid'", "string", fmt.Sprintf("%T", uid)))
	}

	return uidS
}

func apiRefund(t *testing.T, api contracts.Api, request vo.RefundRequest, codeExp int, statusExp string) string {
	resp, err := api.Refund(context.Background(), request)
	if err != nil {
		t.Fatal(sprintfExpAct("api.Authorization: return non nil error", nil, err))
	}
	if resp.StatusCode != codeExp {
		t.Fatal(sprintfExpAct("Response.StatusCode: unexpected status codeExp", codeExp, resp.StatusCode))
	}

	transaction, err := getTransaction(resp.Body)
	if err != nil {
		t.Fatal(sprintfExpAct("getTransaction: return non nil error", nil, err))
	}

	status, ok := transaction["status"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'status'", "key 'status' present", "no such key"))
	}

	statusS, ok := status.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'status'", "string", fmt.Sprintf("%T", status)))
	}

	if statusS != statusExp {
		t.Fatal(sprintfExpAct("transaction: status not equal", statusExp, statusS))
	}

	uid, ok := transaction["uid"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'uid'", "key 'uid' present", "no such key"))
	}

	uidS, ok := uid.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'uid'", "string", fmt.Sprintf("%T", uid)))
	}

	return uidS
}

func apiVoid(t *testing.T, api contracts.Api, request vo.VoidRequest, codeExp int, statusExp string) string {
	resp, err := api.Void(context.Background(), request)
	if err != nil {
		t.Fatal(sprintfExpAct("api.Authorization: return non nil error", nil, err))
	}
	if resp.StatusCode != codeExp {
		t.Fatal(sprintfExpAct("Response.StatusCode: unexpected status codeExp", codeExp, resp.StatusCode))
	}

	transaction, err := getTransaction(resp.Body)
	if err != nil {
		t.Fatal(sprintfExpAct("getTransaction: return non nil error", nil, err))
	}

	status, ok := transaction["status"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'status'", "key 'status' present", "no such key"))
	}

	statusS, ok := status.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'status'", "string", fmt.Sprintf("%T", status)))
	}

	if statusS != statusExp {
		t.Fatal(sprintfExpAct("transaction: status not equal", statusExp, statusS))
	}

	uid, ok := transaction["uid"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'uid'", "key 'uid' present", "no such key"))
	}

	uidS, ok := uid.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'uid'", "string", fmt.Sprintf("%T", uid)))
	}

	return uidS
}

func apiStatusByUid(t *testing.T, api contracts.Api, parentUid string, codeExp int, statusExp string) string {
	resp, err := api.StatusByUid(context.Background(), parentUid)
	if err != nil {
		t.Fatal(sprintfExpAct("api.Authorization: return non nil error", nil, err))
	}
	if resp.StatusCode != codeExp {
		t.Fatal(sprintfExpAct("Response.StatusCode: unexpected status codeExp", codeExp, resp.StatusCode))
	}

	transaction, err := getTransaction(resp.Body)
	if err != nil {
		t.Fatal(sprintfExpAct("getTransaction: return non nil error", nil, err))
	}

	status, ok := transaction["status"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'status'", "key 'status' present", "no such key"))
	}

	statusS, ok := status.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'status'", "string", fmt.Sprintf("%T", status)))
	}

	if statusS != statusExp {
		t.Fatal(sprintfExpAct("transaction: status not equal", statusExp, statusS))
	}

	uid, ok := transaction["uid"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'uid'", "key 'uid' present", "no such key"))
	}

	uidS, ok := uid.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'uid'", "string", fmt.Sprintf("%T", uid)))
	}

	return uidS
}

func apiStatusByTrackingId(t *testing.T, api contracts.Api, trackingid string, codeExp int, statusExp string) string {
	resp, err := api.StatusByTrackingId(context.Background(), trackingid)
	if err != nil {
		t.Fatal(sprintfExpAct("api.Authorization: return non nil error", nil, err))
	}
	if resp.StatusCode != codeExp {
		t.Fatal(sprintfExpAct("Response.StatusCode: unexpected status codeExp", codeExp, resp.StatusCode))
	}

	transactionArray, err := getTransactionsArray(resp.Body)
	if err != nil {
		t.Fatal(sprintfExpAct("getTransaction: return non nil error", nil, err))
	}

	if len(transactionArray) == 0 {
		t.Fatal(sprintfExpAct("transactions: array is empty", "array with at least one element", "empty array"))
	}

	transaction, ok := transactionArray[0].(map[string]interface{})
	if !ok {
		t.Fatal(sprintfExpAct("transactions: value at index 0 isn't assertable as map[string]interface{}", "asserted", "not asserted"))
	}

	status, ok := transaction["status"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'status'", "key 'status' present", "no such key"))
	}

	statusS, ok := status.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'status'", "string", fmt.Sprintf("%T", status)))
	}

	if statusS != statusExp {
		t.Fatal(sprintfExpAct("transaction: status not equal", statusExp, statusS))
	}

	uid, ok := transaction["uid"]
	if !ok {
		t.Fatal(sprintfExpAct("transaction: map doesn't contain key 'uid'", "key 'uid' present", "no such key"))
	}

	uidS, ok := uid.(string)
	if !ok {
		t.Fatal(sprintfExpAct("transaction: not a string value in key 'uid'", "string", fmt.Sprintf("%T", uid)))
	}

	return uidS
}

func getTransaction(body io.Reader) (map[string]interface{}, error) {
	m := map[string]interface{}{}

	err := json.NewDecoder(body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("Decoder.Decode: err is not nil: %w", err)
	}

	if len(m) == 0 {
		return nil, fmt.Errorf("empty response body")
	}

	t, ok := m["transaction"]
	if !ok {
		return nil, fmt.Errorf("no 'transaction' key in map")
	}

	transaction, ok := t.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'transaction' key couldn't be asserted as map[string]interface{}")
	}

	return transaction, nil
}

func getTransactionsArray(body io.ReadCloser) ([]interface{}, error) {
	m := map[string]interface{}{}

	err := json.NewDecoder(body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("Decoder.Decode: err is not nil: %w", err)
	}

	if len(m) == 0 {
		return nil, fmt.Errorf("empty response body")
	}

	t, ok := m["transactions"]
	if !ok {
		return nil, fmt.Errorf("no 'transactions' key in map")
	}

	transaction, ok := t.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'transaction' key couldn't be asserted as []interface{}")
	}

	return transaction, nil
}

func fatalfWithExpectedActual(t *testing.T, msg string, er, ar interface{}) {
	t.Fatalf("%s:\nER: %v,\nAR: %v", msg, er, ar)
}

func sprintfExpAct(msg string, er, ar interface{}) string {
	return fmt.Sprintf("%s:\nExp: %v,\nAct: %v", msg, er, ar)
}

func getUid(body io.Reader) (string, error) {
	m := map[string]interface{}{}

	err := json.NewDecoder(body).Decode(&m)
	if err != nil {
		return "", fmt.Errorf("Decoder.Decode: err is not nil: %w", err)
	}

	if len(m) == 0 {
		return "", fmt.Errorf("response body length == 0")
	}

	transactionMap, ok := m["transaction"]
	if !ok {
		return "", fmt.Errorf("no 'transaction' key in map")
	}

	uid, ok := transactionMap.(map[string]interface{})["uid"]
	if !ok {
		return "", fmt.Errorf("no 'uid' key in transactionMap")
	}

	uidS, ok := uid.(string)
	if !ok {
		return "", fmt.Errorf("value in 'uid' key is not a string")
	}

	return uidS, nil
}

func readBody(t *testing.T, response *http.Response) *bytes.Buffer {
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, response.Body)
	checkError(t, err)
	return buf
}

func logResponse(t *testing.T, response *http.Response) {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		fatalfWithExpectedActual(t, "io.ReadAll error", nil, err)
	}
	t.Logf("Headers: %v", response.Header)
	t.Logf("StatusCode: %d", response.StatusCode)
	t.Logf("Status: %s", response.Status)
	t.Logf("Body: %s", string(b))

}

func checkError(t *testing.T, err error) {
	if err != nil {
		fatalfWithExpectedActual(t, "err is not nil", nil, err)
	}
}
