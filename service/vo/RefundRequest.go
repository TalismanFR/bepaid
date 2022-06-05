package vo

import "encoding/json"

type RefundRequest struct {
	request struct {

		//UID транзакции оплаты или списания средств
		parentUid string `json:"parent_uid"`

		//сумма списания в минимальных денежных единицах, например 1000 для $10.00
		amount Amount `json:"amount"`

		//причина возврата. Максимальная длина: 255 символов
		reason string `json:"reason"`

		//(необязательный) true или false. Параметр управляет процессом проверки входящего запроса на уникальность.
		//Если в течение 30 секунд придет запрос на списание средств с одинаковыми amount и parent_uid, то запрос будет отклонен.
		//По умолчанию, этот параметр имеет значение true
		duplicateCheck *bool `json:"duplicate_check,omitempty"`
	} `json:"request"`
}

func (rr *RefundRequest) ParentUid() string {
	return rr.request.parentUid
}

func (rr *RefundRequest) Amount() Amount {
	return rr.request.amount
}

func (rr *RefundRequest) Reason() string {
	return rr.request.reason
}

func (rr *RefundRequest) DuplicateCheck() *bool {
	return rr.request.duplicateCheck
}

func (rr RefundRequest) MarshalJSON() ([]byte, error) {
	type R struct {
		Request struct {
			ParentUid      string `json:"parent_uid"`
			Amount         Amount `json:"amount"`
			Reason         string `json:"reason"`
			DuplicateCheck *bool  `json:"duplicate_check,omitempty"`
		} `json:"request"`
	}

	r := R{}
	r.Request.ParentUid = rr.ParentUid()
	r.Request.Amount = rr.Amount()
	r.Request.Reason = rr.Reason()
	r.Request.DuplicateCheck = rr.DuplicateCheck()

	return json.Marshal(r)
}

// NewRefundRequest creates RefundRequest with mandatory fields
func NewRefundRequest(parentUid string, amount Amount, reason string) *RefundRequest {
	r := &RefundRequest{}

	r.request.parentUid = parentUid
	r.request.amount = amount
	r.request.reason = reason

	return r
}

func (rr *RefundRequest) WithDuplicateCheck(duplicateCheck bool) *RefundRequest {
	rr.request.duplicateCheck = &duplicateCheck
	return rr
}
