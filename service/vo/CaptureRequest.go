package vo

import (
	"encoding/json"
)

type CaptureRequest struct {

	//UID транзакции авторизации
	parentUid string

	//сумма списания в минимальных денежных единицах, например 1000 для $10.00
	amount Amount

	//(необязательный) true или false. Параметр управляет процессом проверки входящего запроса на уникальность.
	//Если в течение 30 секунд придет запрос на списание средств с одинаковыми amount и parent_uid, то запрос будет отклонен.
	//По умолчанию, этот параметр имеет значение true
	duplicateCheck *bool
}

func (cr *CaptureRequest) ParentUid() string {
	return cr.parentUid
}

func (cr *CaptureRequest) Amount() Amount {
	return cr.amount
}

func (cr *CaptureRequest) DuplicateCheck() *bool {
	return cr.duplicateCheck
}

func (cr CaptureRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Request interface{} `json:"request"`
	}{struct {
		ParentUid      string `json:"parent_uid"`
		Amount         Amount `json:"amount"`
		DuplicateCheck *bool  `json:"duplicate_check,omitempty"`
	}{cr.parentUid, cr.amount, cr.duplicateCheck}})
}

// NewCaptureRequest creates CaptureRequest with mandatory fields
func NewCaptureRequest(parentUid string, amount Amount) *CaptureRequest {
	cr := &CaptureRequest{}

	cr.parentUid = parentUid
	cr.amount = amount

	return cr
}

func (cr *CaptureRequest) WithDuplicateCheck(duplicateCheck bool) *CaptureRequest {
	cr.duplicateCheck = &duplicateCheck
	return cr
}
