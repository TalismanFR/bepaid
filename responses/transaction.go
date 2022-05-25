package responses

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
}

func (tx *TransactionResponse) IsSuccess() bool {
	return tx.Status == success
}
func (tx *TransactionResponse) IsFailed() bool {
	return tx.Status == failed
}

func (tx *TransactionResponse) IsIncomplete() bool {
	return tx.Status == incomplete
}

func (tx *TransactionResponse) IsExpired() bool {
	return tx.Status == expired
}

func (tx *TransactionResponse) IsVoid() bool {
	return tx.Type == void
}

func (tx *TransactionResponse) IsAuthorization() bool {
	return tx.Type == authorization
}

func (tx *TransactionResponse) IsCapture() bool {
	return tx.Type == capture
}

func (tx *TransactionResponse) IsRefund() bool {
	return tx.Type == refund
}

func (tx *TransactionResponse) IsPayment() bool {
	return tx.Type == payment
}
