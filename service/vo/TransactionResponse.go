package vo

const (
	// Status
	success    = "successful"
	failed     = "failed"
	incomplete = "incomplete"
	expired    = "expired"

	// Type
	capture       = "capture"
	authorization = "authorization"
	void          = "void"
	payment       = "payment"
	refund        = "refund"
)

type TransactionResponse struct {
	Transaction struct {
		Message            string `json:"message"`
		RefId              string `json:"ref_id"`
		GatewayId          int    `json:"gateway_id"`
		Uid                string `json:"uid"`
		Status             string `json:"status"`
		MessageTransaction string `json:"message_transaction"`
		Amount             int    `json:"amount"`
		ParentUid          string `json:"parent_uid"`
		ReceiptUrl         string `json:"receipt_url"`
		Currency           string `json:"currency"`
		Type               string `json:"type"`
		Test               bool   `json:"test"`
	} `json:"transaction"`

	// for errors
	Response struct {
		Message string                 `json:"message"`
		Errors  map[string]interface{} `json:"errors"`
	} `json:"response"`
}

func (tr *TransactionResponse) Uid() string {
	return tr.Transaction.Uid
}

func (tr *TransactionResponse) IsSuccess() bool {
	return tr.Transaction.Status == success
}

func (tr *TransactionResponse) IsFailed() bool {
	return tr.Transaction.Status == failed
}

func (tr *TransactionResponse) IsIncomplete() bool {
	return tr.Transaction.Status == incomplete
}

func (tr *TransactionResponse) IsExpired() bool {
	return tr.Transaction.Status == expired
}

func (tr *TransactionResponse) IsVoid() bool {
	return tr.Transaction.Type == void
}

func (tr *TransactionResponse) IsAuthorization() bool {
	return tr.Transaction.Type == authorization
}

func (tr *TransactionResponse) IsCapture() bool {
	return tr.Transaction.Type == capture
}

func (tr *TransactionResponse) IsRefund() bool {
	return tr.Transaction.Type == refund
}

func (tr *TransactionResponse) IsPayment() bool {
	return tr.Transaction.Type == payment
}

func (tr *TransactionResponse) IsError() bool {
	return tr.Response.Message != ""
}

//todo
//методы информации о платеже isSuccess isFailed, isCapture, isVoid, isAuthorization, isRefund, need3ds, expDate time
