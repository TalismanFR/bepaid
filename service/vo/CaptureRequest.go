package vo

import (
	"encoding/json"
)

type CaptureRequest struct {
	request struct {

		//UID транзакции авторизации
		parentUid string `json:"parent_uid"`

		//сумма списания в минимальных денежных единицах, например 1000 для $10.00
		amount Amount `json:"amount"`

		//(необязательный) true или false. Параметр управляет процессом проверки входящего запроса на уникальность.
		//Если в течение 30 секунд придет запрос на списание средств с одинаковыми amount и parent_uid, то запрос будет отклонен.
		//По умолчанию, этот параметр имеет значение true
		duplicateCheck *bool `json:"duplicate_check,omitempty"`
	} `json:"request"`
}

func (cr *CaptureRequest) ParentUid() string {
	return cr.request.parentUid
}

func (cr *CaptureRequest) Amount() Amount {
	return cr.request.amount
}

func (cr *CaptureRequest) DuplicateCheck() *bool {
	return cr.request.duplicateCheck
}

func (cr CaptureRequest) MarshalJSON() ([]byte, error) {
	type C struct {
		Request struct {
			ParentUid      string `json:"parent_uid"`
			Amount         Amount `json:"amount"`
			DuplicateCheck *bool  `json:"duplicate_check,omitempty"`
		} `json:"request"`
	}

	c := C{}
	c.Request.ParentUid = cr.ParentUid()
	c.Request.Amount = cr.Amount()
	c.Request.DuplicateCheck = cr.DuplicateCheck()

	return json.Marshal(c)
}

// NewCaptureRequest creates CaptureRequest with mandatory fields
func NewCaptureRequest(parentUid string, amount Amount) *CaptureRequest {
	r := &CaptureRequest{}

	r.request.parentUid = parentUid
	r.request.amount = amount

	return r
}

func (cr *CaptureRequest) WithDuplicateCheck(duplicateCheck bool) *CaptureRequest {
	cr.request.duplicateCheck = &duplicateCheck
	return cr
}
