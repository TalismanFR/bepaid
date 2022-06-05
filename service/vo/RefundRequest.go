package vo

type RefundRequest struct {
	Request struct {

		//UID транзакции оплаты или списания средств
		ParentUid string `json:"parent_uid"`

		//сумма списания в минимальных денежных единицах, например 1000 для $10.00
		Amount Amount `json:"amount"`

		//причина возврата. Максимальная длина: 255 символов
		Reason string `json:"reason"`

		//(необязательный) true или false. Параметр управляет процессом проверки входящего запроса на уникальность.
		//Если в течение 30 секунд придет запрос на списание средств с одинаковыми amount и parent_uid, то запрос будет отклонен.
		//По умолчанию, этот параметр имеет значение true
		DuplicateCheck *bool `json:"duplicate_check,omitempty"`
	} `json:"request"`
}

// NewRefundRequest creates RefundRequest with mandatory fields
func NewRefundRequest(parentUid string, amount Amount, reason string) *RefundRequest {
	r := &RefundRequest{}

	r.Request.ParentUid = parentUid
	r.Request.Amount = amount
	r.Request.Reason = reason

	return r
}

func (cr *RefundRequest) WithDuplicateCheck(duplicateCheck bool) *RefundRequest {
	cr.Request.DuplicateCheck = &duplicateCheck
	return cr
}
