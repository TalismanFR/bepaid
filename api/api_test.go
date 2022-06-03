package api

import (
	"bepaid-sdk/service/vo"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
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
	// tests, using mockApi, read Request.Body from this channel. So, no parallel testing.
	ch = make(chan io.ReadCloser, 1)

	// testingClient sends Request.Body to ch.
	testingClient = &http.Client{Transport: customRoundTripper{}}

	// mockApi sends Request.Body to ch.
	// Use to check Request.Body marshalling
	mockApi = Api{client: testingClient}

	// correctApi sends requests to real host and return host's response
	// Host address, shop_id, secret retrieved from .env file in 'bepaid' directory
	correctApi *Api
)

type customRoundTripper struct{}

func (customRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	ch <- request.Body
	return nil, nil
}

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		panic("cannot load .env file")
	}

	correctApi = NewApi(http.DefaultClient, os.Getenv("API_HOST"), os.Getenv("SHOP_ID"), os.Getenv("SECRET"))

	m.Run()
}

//////////////////////////////////////
//		Marshalling request body	//
//////////////////////////////////////

func TestApi_PaymentsMarshalRequest(t *testing.T) {

	tests := []struct {
		name string
		req  P
		er   string
	}{
		{"defaultValue", P{}, `{"request":{"amount":0,"currency":"","description":"","tracking_id":"","test":false,"credit_card":{"number":"","verification_value":"","holder":"","exp_month":"","exp_year":"","skip_three_d_secure_verification":false}}}`},
		{"requestConstructor", *vo.NewPaymentRequest(int64(1), "rub", "rub_1", "id1", true, *vo.NewCreditCard("5555", "123", "tim", "05", "2024")), `{"request":{"amount":1,"currency":"rub","description":"rub_1","tracking_id":"id1","test":true,"credit_card":{"number":"5555","verification_value":"123","holder":"tim","exp_month":"05","exp_year":"2024","skip_three_d_secure_verification":false}}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				tc.er,
				func() (*http.Response, error) {
					return mockApi.Payment(context.TODO(), tc.req)
				})
		})
	}
}

func TestApi_AuthorizationsMarshalRequest(t *testing.T) {

	tests := []struct {
		name string
		req  A
		er   string
	}{
		{"defaultValue", A{}, `{"request":{"amount":0,"currency":"","description":"","tracking_id":"","test":false,"credit_card":{"number":"","verification_value":"","holder":"","exp_month":"","exp_year":"","skip_three_d_secure_verification":false}}}`},
		{"requestConstructor", *vo.NewAuthorizationRequest(int64(1), "rub", "rub_1", "id1", true, *vo.NewCreditCard("5555", "123", "tim", "05", "2024")), `{"request":{"amount":1,"currency":"rub","description":"rub_1","tracking_id":"id1","test":true,"credit_card":{"number":"5555","verification_value":"123","holder":"tim","exp_month":"05","exp_year":"2024","skip_three_d_secure_verification":false}}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				tc.er,
				func() (*http.Response, error) {
					return mockApi.Authorization(context.TODO(), tc.req)
				})
		})
	}
}

func TestApi_CapturesMarshalRequest(t *testing.T) {

	tests := []struct {
		name string
		req  vo.CaptureRequest
		er   string
	}{
		{"defaultValue", C{}, `{"request":{"parent_uid":"","amount":0}}`},
		{"requestConstructor", *vo.NewCaptureRequest("id123", int64(63)), `{"request":{"parent_uid":"id123","amount":63}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				tc.er,
				func() (*http.Response, error) {
					return mockApi.Capture(context.TODO(), tc.req)
				})
		})
	}
}

func TestApi_VoidsMarshalRequest(t *testing.T) {

	tests := []struct {
		name string
		req  V
		er   string
	}{
		{"defaultValue", V{}, `{"request":{"parent_uid":"","amount":0}}`},
		{"requestConstructor", *vo.NewVoidRequest("id123", int64(63)), `{"request":{"parent_uid":"id123","amount":63}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				tc.er,
				func() (*http.Response, error) {
					return mockApi.Void(context.TODO(), tc.req)
				})
		})
	}
}

func TestApi_RefundsMarshalRequest(t *testing.T) {

	tests := []struct {
		name string
		req  R
		er   string
	}{
		{"defaultValue", R{}, `{"request":{"parent_uid":"","amount":0,"reason":""}}`},
		{"requestConstructor", *vo.NewRefundRequest("id123", int64(63), "reason"), `{"request":{"parent_uid":"id123","amount":63,"reason":"reason"}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testMarshallRequest(
				t,
				tc.er,
				func() (*http.Response, error) {
					return mockApi.Refund(context.TODO(), tc.req)
				})
		})
	}
}

//////////////////////////////////////
// 			Single requests	 		//
//////////////////////////////////////

func TestApi_Authorization(t *testing.T) {
	tests := []struct {
		name   string
		a      A
		code   int
		status string
	}{
		{"test1",
			*vo.NewAuthorizationRequest(101, "RUB", "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
			http.StatusOK,
			"successful",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("UID: %s", apiAuthorization(t, tc.a, tc.code, tc.status))
		})
	}
}

func TestApi_Payment(t *testing.T) {

	tests := []struct {
		name   string
		p      P
		code   int
		status string
	}{
		{"test1",
			*vo.NewPaymentRequest(int64(100), "RUB", "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")).WithDuplicateCheck(false),
			http.StatusOK,
			"successful",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("UID: %s", apiPayment(t, tc.p, tc.code, tc.status))
		})
	}
}

func TestApi_StatusByTrackingId(t *testing.T) {
	r, _ := correctApi.StatusByTrackingId(context.Background(), "mytrackingid")

	defer r.Body.Close()
	buf := bytes.Buffer{}
	b, _ := io.ReadAll(r.Body)
	json.Indent(&buf, b, "", "\t")
	fmt.Println(string(buf.Bytes()))
}

func TestApi_StatusByUid(t *testing.T) {
	r, _ := correctApi.StatusByUid(context.Background(), "151534003-9d0e9c9aa1")

	defer r.Body.Close()
	buf := bytes.Buffer{}
	b, _ := io.ReadAll(r.Body)
	json.Indent(&buf, b, "", "\t")
	fmt.Println(string(buf.Bytes()))
}

//////////////////////////////////////
//		Sequential requests			//
//////////////////////////////////////

func TestApi_AuthorizationCapture(t *testing.T) {

	type (
		Auth struct {
			a      A
			code   int
			status string
		}
		Capt struct {
			c      C
			code   int
			status string
		}
	)
	amount := rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100

	tests := []struct {
		name string
		auth Auth
		capt Capt
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, "RUB", "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			Capt{
				*vo.NewCaptureRequest("", amount),
				http.StatusOK,
				"successful"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			uid := apiAuthorization(t, tc.auth.a, tc.auth.code, tc.auth.status)
			t.Logf("A.Uid: %s", uid)

			tc.capt.c.Request.ParentUid = uid

			uid = apiCapture(t, tc.capt.c, tc.capt.code, tc.capt.status)
			t.Logf("C.Uid: %s", uid)
		})
	}

}

func TestApi_AuthorizationVoid(t *testing.T) {
	type (
		Auth struct {
			a      A
			code   int
			status string
		}
		Void struct {
			v      V
			code   int
			status string
		}
	)
	amount := rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100

	tests := []struct {
		name string
		auth Auth
		void Void
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, "RUB", "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			Void{
				*vo.NewVoidRequest("", amount),
				http.StatusOK,
				"successful"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			uid := apiAuthorization(t, tc.auth.a, tc.auth.code, tc.auth.status)
			t.Logf("A.Uid: %s", uid)

			tc.void.v.Request.ParentUid = uid

			uid = apiVoid(t, tc.void.v, tc.void.code, tc.void.status)
			t.Logf("V.Uid: %s", uid)
		})
	}
}

func TestApi_AuthorizationCaptureRefund(t *testing.T) {
	type (
		Auth struct {
			a      A
			code   int
			status string
		}
		Capt struct {
			c      C
			code   int
			status string
		}
		Refund struct {
			r      R
			code   int
			status string
		}
	)
	amount := rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100

	tests := []struct {
		name   string
		auth   Auth
		capt   Capt
		refund Refund
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, "RUB", "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			Capt{
				*vo.NewCaptureRequest("", amount),
				http.StatusOK,
				"successful"},
			Refund{
				*vo.NewRefundRequest("", amount, "give me my money back"),
				http.StatusOK,
				"successful",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			uid := apiAuthorization(t, tc.auth.a, tc.auth.code, tc.auth.status)
			t.Logf("A.Uid: %s", uid)

			tc.capt.c.Request.ParentUid = uid

			uid = apiCapture(t, tc.capt.c, tc.capt.code, tc.capt.status)
			t.Logf("C.Uid: %s", uid)

			tc.refund.r.Request.ParentUid = uid

			uid = apiRefund(t, tc.refund.r, tc.refund.code, tc.refund.status)
			t.Logf("R.Uid: %s", uid)
		})
	}
}

func TestApi_PaymentRefund(t *testing.T) {
	type (
		Paym struct {
			p      P
			code   int
			status string
		}
		Refund struct {
			r      R
			code   int
			status string
		}
	)
	amount := rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100

	tests := []struct {
		name   string
		paym   Paym
		refund Refund
	}{
		{"positiveTest1",
			Paym{
				*vo.NewPaymentRequest(amount, "RUB", "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
				http.StatusOK,
				"successful"},
			Refund{
				*vo.NewRefundRequest("", amount, "give me my money back"),
				http.StatusOK,
				"successful"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			uid := apiPayment(t, tc.paym.p, tc.paym.code, tc.paym.status)
			t.Logf("A.Uid: %s", uid)

			tc.refund.r.Request.ParentUid = uid

			uid = apiRefund(t, tc.refund.r, tc.refund.code, tc.refund.status)
			t.Logf("R.Uid: %s", uid)
		})
	}
}

func TestApi_AuthorizationStatusByUid(t *testing.T) {
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
	amount := rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100

	tests := []struct {
		name    string
		auth    Auth
		statUid StatusUid
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, "RUB", "it's description", "mytrackingid", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
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
		t.Run(tc.name, func(t *testing.T) {

			uid := apiAuthorization(t, tc.auth.a, tc.auth.code, tc.auth.status)
			t.Logf("A.Uid: %s", uid)

			tc.statUid.uid = uid

			uid = apiStatusByUid(t, tc.statUid.uid, tc.statUid.code, tc.statUid.status)
			t.Logf("S.Uid: %s", uid)
		})
	}
}

func TestApi_AuthorizationStatusByTrackingId(t *testing.T) {
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
	amount := rand.New(rand.NewSource(time.Now().Unix())).Int63() % 100

	tests := []struct {
		name       string
		auth       Auth
		statusTrId StatusTrId
	}{
		{"positiveTest1",
			Auth{
				*vo.NewAuthorizationRequest(amount, "RUB", "it's description", "mytrackingid5678", true, *vo.NewCreditCard("4200000000000000", "123", "tim", "01", "2024")),
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
		t.Run(tc.name, func(t *testing.T) {

			uid := apiAuthorization(t, tc.auth.a, tc.auth.code, tc.auth.status)
			t.Logf("A.Uid: %s", uid)

			uid = apiStatusByTrackingId(t, tc.statusTrId.trackingId, tc.statusTrId.code, tc.statusTrId.status)
			t.Logf("S.Uid: %s", uid)
		})
	}
}

//////////////////////////////////////
//		Helper functions			//
//////////////////////////////////////

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

func apiPayment(t *testing.T, request vo.PaymentRequest, codeExp int, statusExp string) string {
	resp, err := correctApi.Payment(context.Background(), request)
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

func apiAuthorization(t *testing.T, request vo.AuthorizationRequest, codeExp int, statusExp string) string {
	resp, err := correctApi.Authorization(context.Background(), request)
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

func apiCapture(t *testing.T, request vo.CaptureRequest, codeExp int, statusExp string) string {
	resp, err := correctApi.Capture(context.Background(), request)
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

func apiRefund(t *testing.T, request vo.RefundRequest, codeExp int, statusExp string) string {
	resp, err := correctApi.Refund(context.Background(), request)
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

func apiVoid(t *testing.T, request vo.VoidRequest, codeExp int, statusExp string) string {
	resp, err := correctApi.Void(context.Background(), request)
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

func apiStatusByUid(t *testing.T, parentUid string, codeExp int, statusExp string) string {
	resp, err := correctApi.StatusByUid(context.Background(), parentUid)
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

func apiStatusByTrackingId(t *testing.T, trackingid string, codeExp int, statusExp string) string {
	resp, err := correctApi.StatusByTrackingId(context.Background(), trackingid)
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

func testMarshallRequest(t *testing.T, er string, startRequest func() (*http.Response, error)) {
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
