package vo

import (
	"encoding/json"
	"time"
)

type PaymentRequest struct {
	request struct {

		//стоимость в минимальных денежных единицах.
		//Например, $32.45 должна быть отправлена как 3245
		amount Amount `json:"amount"`

		//валюта в ISO-4217 формате, например USD
		currency Currency `json:"currency"`

		//описание заказа. Максимальная длина: 255 символов
		description string `json:"description"`

		//id транзакции или заказа в вашей системе.
		//Максимальная длина: 255 символов.
		//Пожалуйста, используйте уникальное значение для того, чтобы при запросе статуса транзакции получить актуальную информацию.
		//В противном случае вы получите первую найденную по tracking_id транзакцию
		trackingId string `json:"tracking_id"`

		//(необязательно) время в формате ISO 8601, до которого должна быть завершена операция.
		//По умолчанию - бессрочно.
		//Формат: YYYY-MM-DDThh:mm:ssTZD, где YYYY – год (например 2019), MM – месяц (например 02), DD – день (например 09), hh – часы (например 18), mm – минуты (например 20), ss – секунды (например 45), TZD – часовой пояс (+hh:mm или –hh:mm), например +03:00 для Минска.
		//Если в указанный момент платёж всё ещё не будет оплачен, он будет переведён в статус expired
		expiredAt *time.Time `json:"expired_at,omitempty"`

		//(необязательный) true или false.
		//Параметр управляет процессом проверки входящего запроса на уникальность.
		//Если в течение 30 секунд придет запрос на авторизацию с одинаковыми amount и number или token, то запрос будет отклонен.
		//По умолчанию, этот параметр имеет значение true
		duplicateCheck *bool `json:"duplicate_check,omitempty"`

		//параметр обязателен, если 3-D Secure включен.
		//Обратитесь к менеджеру за информацией. return_url - это URL на стороне торговца,
		//на который bePaid будет перенаправлять клиента после возврата с 3-D Secure проверки
		returnUrl string `json:"return_url,omitempty"`

		//true или false. Транзакция будет тестовой, если значение true.
		test bool `json:"test"`

		creditCard CreditCard `json:"credit_card"`

		//секция, содержащая дополнительную информацию о платеже
		additionalData map[string]interface{} `json:"additional_data,omitempty"`

		customer *Customer `json:"customer,omitempty"`
	} `json:"request"`
}

func (pr PaymentRequest) MarshalJSON() ([]byte, error) {
	type P struct {
		Request struct {
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
		} `json:"request"`
	}

	p := P{}
	p.Request.Amount = pr.Amount()
	p.Request.Currency = pr.Currency()
	p.Request.Description = pr.Description()
	p.Request.TrackingId = pr.TrackingId()
	p.Request.ExpiredAt = pr.ExpiredAt()
	p.Request.DuplicateCheck = pr.DuplicateCheck()
	p.Request.ReturnUrl = pr.ReturnUrl()
	p.Request.Test = pr.Test()
	p.Request.CreditCard = pr.CreditCard()
	p.Request.AdditionalData = pr.AdditionalData()
	p.Request.Customer = pr.Customer()

	return json.Marshal(p)
}

func (pr *PaymentRequest) Amount() Amount {
	return pr.request.amount
}

func (pr *PaymentRequest) Currency() Currency {
	return pr.request.currency
}

func (pr *PaymentRequest) Description() string {
	return pr.request.description
}

func (pr *PaymentRequest) TrackingId() string {
	return pr.request.trackingId
}

func (pr *PaymentRequest) ExpiredAt() *time.Time {
	return pr.request.expiredAt
}

func (pr *PaymentRequest) DuplicateCheck() *bool {
	return pr.request.duplicateCheck
}

func (pr *PaymentRequest) ReturnUrl() string {
	return pr.request.returnUrl
}

func (pr *PaymentRequest) Test() bool {
	return pr.request.test
}

func (pr *PaymentRequest) CreditCard() CreditCard {
	return pr.request.creditCard
}

func (pr *PaymentRequest) AdditionalData() map[string]interface{} {
	return pr.request.additionalData
}

func (pr *PaymentRequest) Customer() *Customer {
	return pr.request.customer
}

// NewPaymentRequest creates PaymentRequest with mandatory fields
func NewPaymentRequest(amount Amount, currency Currency, description string, trackingId string, test bool, cc CreditCard) *PaymentRequest {
	r := &PaymentRequest{}

	r.request.amount = amount
	r.request.currency = currency
	r.request.description = description
	r.request.trackingId = trackingId
	r.request.test = test
	r.request.creditCard = cc

	return r
}

func (pr *PaymentRequest) WithExpiresAt(expiresAt time.Time) *PaymentRequest {
	pr.request.expiredAt = &expiresAt
	return pr
}

func (cr *PaymentRequest) WithDuplicateCheck(duplicateCheck bool) *PaymentRequest {
	cr.request.duplicateCheck = &duplicateCheck
	return cr
}

func (a *PaymentRequest) WithReturnUrl(returnUrl string) *PaymentRequest {
	a.request.returnUrl = returnUrl
	return a
}

// WithAdditionalData saves argument to Authorizationrequest.request.AdditionalData field.
//
// Don't change content of additionalData after function call.
func (a *PaymentRequest) WithAdditionalData(additionalData map[string]interface{}) *PaymentRequest {
	a.request.additionalData = additionalData
	return a
}

func (a *PaymentRequest) WithCustomer(customer Customer) *PaymentRequest {
	a.request.customer = &customer
	return a
}
