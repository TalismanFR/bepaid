package vo

import (
	"encoding/json"
	"time"
)

type PaymentRequest struct {

	//стоимость в минимальных денежных единицах.
	//Например, $32.45 должна быть отправлена как 3245
	amount Amount

	//валюта в ISO-4217 формате, например USD
	currency Currency

	//описание заказа. Максимальная длина: 255 символов
	description string

	//id транзакции или заказа в вашей системе.
	//Максимальная длина: 255 символов.
	//Пожалуйста, используйте уникальное значение для того, чтобы при запросе статуса транзакции получить актуальную информацию.
	//В противном случае вы получите первую найденную по tracking_id транзакцию
	trackingId string

	//(необязательно) время в формате ISO 8601, до которого должна быть завершена операция.
	//По умолчанию - бессрочно.
	//Формат: YYYY-MM-DDThh:mm:ssTZD, где YYYY – год (например 2019), MM – месяц (например 02), DD – день (например 09), hh – часы (например 18), mm – минуты (например 20), ss – секунды (например 45), TZD – часовой пояс (+hh:mm или –hh:mm), например +03:00 для Минска.
	//Если в указанный момент платёж всё ещё не будет оплачен, он будет переведён в статус expired
	expiredAt *time.Time

	//(необязательный) true или false.
	//Параметр управляет процессом проверки входящего запроса на уникальность.
	//Если в течение 30 секунд придет запрос на авторизацию с одинаковыми amount и number или token, то запрос будет отклонен.
	//По умолчанию, этот параметр имеет значение true
	duplicateCheck *bool

	//параметр обязателен, если 3-D Secure включен.
	//Обратитесь к менеджеру за информацией. return_url - это URL на стороне торговца,
	//на который bePaid будет перенаправлять клиента после возврата с 3-D Secure проверки
	returnUrl string

	//true или false. Транзакция будет тестовой, если значение true.
	test bool

	creditCard CreditCard

	//секция, содержащая дополнительную информацию о платеже
	additionalData map[string]interface{}

	customer *Customer
}

func (pr PaymentRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Request interface{} `json:"request"`
	}{struct {
		Amount         Amount                 `json:"amount"`
		Currency       Currency               `json:"currency"`
		Description    string                 `json:"description"`
		TrackingId     string                 `json:"tracking_id"`
		ExpiredAt      *time.Time             `json:"expired_at,omitempty"`
		DuplicateCheck *bool                  `json:"duplicate_check,omitempty"`
		ReturnUrl      string                 `json:"return_url,omitempty"`
		Test           bool                   `json:"test"`
		CreditCard     CreditCard             `json:"credit_card"`
		AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
		Customer       *Customer              `json:"customer,omitempty"`
	}{
		pr.amount,
		pr.currency,
		pr.description,
		pr.trackingId,
		pr.expiredAt,
		pr.duplicateCheck,
		pr.returnUrl,
		pr.test,
		pr.creditCard,
		pr.additionalData,
		pr.customer,
	}})
}

func (pr *PaymentRequest) Amount() Amount {
	return pr.amount
}

func (pr *PaymentRequest) Currency() Currency {
	return pr.currency
}

func (pr *PaymentRequest) Description() string {
	return pr.description
}

func (pr *PaymentRequest) TrackingId() string {
	return pr.trackingId
}

func (pr *PaymentRequest) ExpiredAt() *time.Time {
	return pr.expiredAt
}

func (pr *PaymentRequest) DuplicateCheck() *bool {
	return pr.duplicateCheck
}

func (pr *PaymentRequest) ReturnUrl() string {
	return pr.returnUrl
}

func (pr *PaymentRequest) Test() bool {
	return pr.test
}

func (pr *PaymentRequest) CreditCard() CreditCard {
	return pr.creditCard
}

func (pr *PaymentRequest) AdditionalData() map[string]interface{} {
	return pr.additionalData
}

func (pr *PaymentRequest) Customer() *Customer {
	return pr.customer
}

// NewPaymentRequest creates PaymentRequest with mandatory fields
func NewPaymentRequest(amount Amount, currency Currency, description string, trackingId string, test bool, cc CreditCard) *PaymentRequest {
	r := &PaymentRequest{}

	r.amount = amount
	r.currency = currency
	r.description = description
	r.trackingId = trackingId
	r.test = test
	r.creditCard = cc

	return r
}

func (pr *PaymentRequest) WithExpiresAt(expiresAt time.Time) *PaymentRequest {
	pr.expiredAt = &expiresAt
	return pr
}

func (pr *PaymentRequest) WithDuplicateCheck(duplicateCheck bool) *PaymentRequest {
	pr.duplicateCheck = &duplicateCheck
	return pr
}

func (pr *PaymentRequest) WithReturnUrl(returnUrl string) *PaymentRequest {
	pr.returnUrl = returnUrl
	return pr
}

// WithAdditionalData saves argument to Authorizationrequest.AdditionalData field.
//
// Don't change content of additionalData after function call.
func (pr *PaymentRequest) WithAdditionalData(additionalData map[string]interface{}) *PaymentRequest {
	pr.additionalData = additionalData
	return pr
}

func (pr *PaymentRequest) WithCustomer(customer Customer) *PaymentRequest {
	pr.customer = &customer
	return pr
}
