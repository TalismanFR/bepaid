package vo

import (
	"time"
)

type PaymentRequest struct {
	Request struct {

		//стоимость в минимальных денежных единицах.
		//Например, $32.45 должна быть отправлена как 3245
		Amount Amount `json:"amount"`

		//валюта в ISO-4217 формате, например USD
		Currency string `json:"currency"`

		//описание заказа. Максимальная длина: 255 символов
		Description string `json:"description"`

		//id транзакции или заказа в вашей системе.
		//Максимальная длина: 255 символов.
		//Пожалуйста, используйте уникальное значение для того, чтобы при запросе статуса транзакции получить актуальную информацию.
		//В противном случае вы получите первую найденную по tracking_id транзакцию
		TrackingId string `json:"tracking_id"`

		//(необязательно) время в формате ISO 8601, до которого должна быть завершена операция.
		//По умолчанию - бессрочно.
		//Формат: YYYY-MM-DDThh:mm:ssTZD, где YYYY – год (например 2019), MM – месяц (например 02), DD – день (например 09), hh – часы (например 18), mm – минуты (например 20), ss – секунды (например 45), TZD – часовой пояс (+hh:mm или –hh:mm), например +03:00 для Минска.
		//Если в указанный момент платёж всё ещё не будет оплачен, он будет переведён в статус expired
		ExpiredAt *time.Time `json:"expired_at,omitempty"`

		//(необязательный) true или false.
		//Параметр управляет процессом проверки входящего запроса на уникальность.
		//Если в течение 30 секунд придет запрос на авторизацию с одинаковыми amount и number или token, то запрос будет отклонен.
		//По умолчанию, этот параметр имеет значение true
		DuplicateCheck *bool `json:"duplicate_check,omitempty"`

		//параметр обязателен, если 3-D Secure включен.
		//Обратитесь к менеджеру за информацией.
		//return_url - это URL на стороне торговца, на который
		//bePaid будет перенаправлять клиента после возврата с 3-D Secure проверки
		ReturnUrl string `json:"return_url,omitempty"`

		//true или false. Транзакция будет тестовой, если значение true.
		Test bool `json:"test"`

		CreditCard CreditCard `json:"credit_card"`

		//секция, содержащая дополнительную информацию о платеже
		AdditionalData map[string]interface{} `json:"additional_data,omitempty"`

		Customer *Customer `json:"customer,omitempty"`
	} `json:"request"`
}

// NewPaymentRequest creates PaymentRequest with mandatory fields
func NewPaymentRequest(amount Amount, currency, description, trackingId string, test bool, cc CreditCard) *PaymentRequest {
	r := &PaymentRequest{}

	r.Request.Amount = amount
	r.Request.Currency = currency
	r.Request.Description = description
	r.Request.TrackingId = trackingId
	r.Request.Test = test
	r.Request.CreditCard = cc

	return r
}

func (pr *PaymentRequest) WithExpiresAt(expiresAt time.Time) *PaymentRequest {
	pr.Request.ExpiredAt = &expiresAt
	return pr
}

func (cr *PaymentRequest) WithDuplicateCheck(duplicateCheck bool) *PaymentRequest {
	cr.Request.DuplicateCheck = &duplicateCheck
	return cr
}

func (a *PaymentRequest) WithReturnUrl(returnUrl string) *PaymentRequest {
	a.Request.ReturnUrl = returnUrl
	return a
}

// WithAdditionalData saves argument to AuthorizationRequest.Request.AdditionalData field.
//
// Don't change content of additionalData after function call.
func (a *PaymentRequest) WithAdditionalData(additionalData map[string]interface{}) *PaymentRequest {
	a.Request.AdditionalData = additionalData
	return a
}

func (a *PaymentRequest) WithCustomer(customer Customer) *PaymentRequest {
	a.Request.Customer = &customer
	return a
}
