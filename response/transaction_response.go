package response

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
		Uid                string `json:"uid"`
		Status             string `json:"status"`
		Amount             int    `json:"amount"`
		Message            string `json:"message"`
		Currency           string `json:"currency"`
		RefId              string `json:"ref_id"`
		GatewayId          int    `json:"gateway_id"`
		MessageTransaction string `json:"message_transaction"`
		ParentUid          string `json:"parent_uid"`
		ReceiptUrl         string `json:"receipt_url"`
		Type               string `json:"type"`
		Test               bool   `json:"test"`
	} `json:"transaction"`
}

func (tx *TransactionResponse) IsSuccess() bool {
	return tx.Transaction.Status == success
}
func (tx *TransactionResponse) IsFailed() bool {
	return tx.Transaction.Status == failed
}

func (tx *TransactionResponse) IsIncomplete() bool {
	return tx.Transaction.Status == incomplete
}

func (tx *TransactionResponse) IsExpired() bool {
	return tx.Transaction.Status == expired
}

func (tx *TransactionResponse) IsVoid() bool {
	return tx.Transaction.Type == void
}

func (tx *TransactionResponse) IsAuthorization() bool {
	return tx.Transaction.Type == authorization
}

func (tx *TransactionResponse) IsCapture() bool {
	return tx.Transaction.Type == capture
}

func (tx *TransactionResponse) IsRefund() bool {
	return tx.Transaction.Type == refund
}

func (tx *TransactionResponse) IsPayment() bool {
	return tx.Transaction.Type == payment
}

//todo
//методы информации о платеже isSuccess isFailed, isCapture, isVoid, isAuthorization, isRefund, need3ds, expDate time
