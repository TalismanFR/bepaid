package vo

type VoidRequest struct {
	Request struct {

		//UID транзакции авторизации
		ParentUid string `json:"parent_uid"`

		//сумма списания в минимальных денежных единицах, например 1000 для $10.00
		Amount Amount `json:"amount"`

		//(необязательный) true или false. Параметр управляет процессом проверки входящего запроса на уникальность.
		//Если в течение 30 секунд придет запрос на списание средств с одинаковыми amount и parent_uid, то запрос будет отклонен.
		//По умолчанию, этот параметр имеет значение true
		DuplicateCheck *bool `json:"duplicate_check,omitempty"`
	} `json:"request"`
}

// NewVoidRequest creates VoidRequest with mandatory fields
func NewVoidRequest(parentUid string, amount Amount) *VoidRequest {
	r := &VoidRequest{}

	r.Request.ParentUid = parentUid
	r.Request.Amount = amount

	return r
}

func (cr *VoidRequest) WithDuplicateCheck(duplicateCheck bool) *VoidRequest {
	cr.Request.DuplicateCheck = &duplicateCheck
	return cr
}
