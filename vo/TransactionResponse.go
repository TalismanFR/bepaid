package vo

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
}

//todo
//методы информации о платеже isSuccess isFailed, isCapture, isVoid, isAuthorization, isRefund, need3ds, expDate time
