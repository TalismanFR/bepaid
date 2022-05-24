package vo

type CaptureRequest struct {
	Request struct {
		//UID транзакции авторизации
		ParentUid string `json:"parent_uid"`
		//сумма списания в минимальных денежных единицах, например 1000 для $10.00
		Amount int64 `json:"amount"`
		//(необязательный) true или false. Параметр управляет процессом проверки входящего запроса на уникальность. Если в течение 30 секунд придет запрос на списание средств с одинаковыми amount и parent_uid, то запрос будет отклонен. По умолчанию, этот параметр имеет значение true
		DuplicateCheck bool `json:"duplicate_check"`
	} `json:"request"`
}

func NewCaptureRequest() *CaptureRequest {
	return &CaptureRequest{}
}
