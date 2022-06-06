package vo

import "encoding/json"

type VoidRequest struct {

	//UID транзакции авторизации
	parentUid string

	//сумма списания в минимальных денежных единицах, например 1000 для $10.00
	amount Amount

	//(необязательный) true или false. Параметр управляет процессом проверки входящего запроса на уникальность.
	//Если в течение 30 секунд придет запрос на списание средств с одинаковыми amount и parent_uid, то запрос будет отклонен.
	//По умолчанию, этот параметр имеет значение true
	duplicateCheck *bool
}

func (vr VoidRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Request interface{} `json:"request"`
	}{struct {
		ParentUid      string `json:"parent_uid"`
		Amount         Amount `json:"amount"`
		DuplicateCheck *bool  `json:"duplicate_check,omitempty"`
	}{
		vr.parentUid,
		vr.amount,
		vr.duplicateCheck,
	}})

}

func (vr *VoidRequest) ParentUid() string {
	return vr.parentUid
}

func (vr *VoidRequest) Amount() Amount {
	return vr.amount
}

func (vr *VoidRequest) DuplicateCheck() *bool {
	return vr.duplicateCheck
}

// NewVoidRequest creates VoidRequest with mandatory fields
func NewVoidRequest(parentUid string, amount Amount) *VoidRequest {
	r := &VoidRequest{}

	r.parentUid = parentUid
	r.amount = amount

	return r
}

func (vr *VoidRequest) WithDuplicateCheck(duplicateCheck bool) *VoidRequest {
	vr.duplicateCheck = &duplicateCheck
	return vr
}
