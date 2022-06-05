package vo

import "encoding/json"

type VoidRequest struct {
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

func (vr VoidRequest) MarshalJSON() ([]byte, error) {
	type V struct {
		Request struct {
			ParentUid      string `json:"parent_uid"`
			Amount         Amount `json:"amount"`
			DuplicateCheck *bool  `json:"duplicate_check,omitempty"`
		} `json:"request"`
	}

	v := V{}
	v.Request.ParentUid = vr.ParentUid()
	v.Request.Amount = vr.Amount()
	v.Request.DuplicateCheck = vr.DuplicateCheck()

	return json.Marshal(v)

}

func (vr *VoidRequest) ParentUid() string {
	return vr.request.parentUid
}

func (vr *VoidRequest) Amount() Amount {
	return vr.request.amount
}

func (vr *VoidRequest) DuplicateCheck() *bool {
	return vr.request.duplicateCheck
}

// NewVoidRequest creates VoidRequest with mandatory fields
func NewVoidRequest(parentUid string, amount Amount) *VoidRequest {
	r := &VoidRequest{}

	r.request.parentUid = parentUid
	r.request.amount = amount

	return r
}

func (vr *VoidRequest) WithDuplicateCheck(duplicateCheck bool) *VoidRequest {
	vr.request.duplicateCheck = &duplicateCheck
	return vr
}
