package vo

import "encoding/json"

type RefundRequest struct {

	//UID транзакции оплаты или списания средств
	parentUid string

	//сумма списания в минимальных денежных единицах, например 1000 для $10.00
	amount Amount

	//причина возврата. Максимальная длина: 255 символов
	reason string

	//(необязательный) true или false. Параметр управляет процессом проверки входящего запроса на уникальность.
	//Если в течение 30 секунд придет запрос на списание средств с одинаковыми amount и parent_uid, то запрос будет отклонен.
	//По умолчанию, этот параметр имеет значение true
	duplicateCheck *bool
}

func (rr *RefundRequest) ParentUid() string {
	return rr.parentUid
}

func (rr *RefundRequest) Amount() Amount {
	return rr.amount
}

func (rr *RefundRequest) Reason() string {
	return rr.reason
}

func (rr *RefundRequest) DuplicateCheck() *bool {
	return rr.duplicateCheck
}

func (rr RefundRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Request interface{} `json:"request"`
	}{struct {
		ParentUid      string `json:"parent_uid"`
		Amount         Amount `json:"amount"`
		Reason         string `json:"reason"`
		DuplicateCheck *bool  `json:"duplicate_check,omitempty"`
	}{rr.parentUid,
		rr.amount,
		rr.reason,
		rr.duplicateCheck}})
}

// NewRefundRequest creates RefundRequest with mandatory fields
func NewRefundRequest(parentUid string, amount Amount, reason string) *RefundRequest {
	r := &RefundRequest{}

	r.parentUid = parentUid
	r.amount = amount
	r.reason = reason

	return r
}

func (rr *RefundRequest) WithDuplicateCheck(duplicateCheck bool) *RefundRequest {
	rr.duplicateCheck = &duplicateCheck
	return rr
}
