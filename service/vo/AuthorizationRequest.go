package vo

import (
	"encoding/json"
	"fmt"
)

type AuthorizationRequest struct {
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

func (ar AuthorizationRequest) MarshalJSON() ([]byte, error) {
	type A struct {
		Request struct {
			Amount         Amount                 `json:"amount"`
			Currency       Currency               `json:"currency"`
			Description    string                 `json:"description"`
			TrackingId     string                 `json:"tracking_id"`
			DuplicateCheck *bool                  `json:"duplicate_check,omitempty"`
			ReturnUrl      string                 `json:"return_url,omitempty"`
			Test           bool                   `json:"test"`
			CreditCard     CreditCard             `json:"credit_card"`
			AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
			Customer       *Customer              `json:"customer,omitempty"`
		} `json:"request"`
	}

	a := A{}
	a.Request.Amount = ar.Amount()
	a.Request.Currency = ar.Currency()
	a.Request.Description = ar.Description()
	a.Request.TrackingId = ar.TrackingId()
	a.Request.DuplicateCheck = ar.DuplicateCheck()
	a.Request.ReturnUrl = ar.ReturnUrl()
	a.Request.Test = ar.Test()
	a.Request.CreditCard = ar.CreditCard()
	a.Request.AdditionalData = ar.AdditionalData()
	a.Request.Customer = ar.Customer()

	return json.Marshal(a)
}

func (ar *AuthorizationRequest) Amount() Amount {
	return ar.request.amount
}

func (ar *AuthorizationRequest) Currency() Currency {
	return ar.request.currency
}

func (ar *AuthorizationRequest) Description() string {
	return ar.request.description
}

func (ar *AuthorizationRequest) TrackingId() string {
	return ar.request.trackingId
}

func (ar *AuthorizationRequest) DuplicateCheck() *bool {
	return ar.request.duplicateCheck
}

func (ar *AuthorizationRequest) ReturnUrl() string {
	return ar.request.returnUrl
}

func (ar *AuthorizationRequest) Test() bool {
	return ar.request.test
}

func (ar *AuthorizationRequest) CreditCard() CreditCard {
	return ar.request.creditCard
}

func (ar *AuthorizationRequest) AdditionalData() map[string]interface{} {
	return ar.request.additionalData
}

func (ar *AuthorizationRequest) Customer() *Customer {
	return ar.request.customer
}

// NewAuthorizationRequest creates AuthorizationRequest with mandatory fields
func NewAuthorizationRequest(amount Amount, currency Currency, description string, trackingId string, test bool, cc CreditCard) *AuthorizationRequest {

	r := &AuthorizationRequest{}

	r.request.amount = amount
	r.request.currency = currency
	r.request.description = description
	r.request.trackingId = trackingId
	r.request.test = test
	r.request.creditCard = cc

	return r
}

func (ar *AuthorizationRequest) Valid() (bool, error) {
	if len(ar.request.description) > 255 {
		return false, fmt.Errorf("description length greater than 255")
	}
	if len(ar.request.trackingId) > 255 {
		return false, fmt.Errorf("trackingId length greater than 255")
	}

	return true, nil
}

func (ar *AuthorizationRequest) WithDuplicateCheck(duplicateCheck bool) *AuthorizationRequest {
	ar.request.duplicateCheck = &duplicateCheck
	return ar
}

func (ar *AuthorizationRequest) WithReturnUrl(returnUrl string) *AuthorizationRequest {
	ar.request.returnUrl = returnUrl
	return ar
}

// WithAdditionalData saves argument to Authorizationrequest.request.AdditionalData field.
//
// Don't change content of additionalData after function call.
func (ar *AuthorizationRequest) WithAdditionalData(additionalData map[string]interface{}) *AuthorizationRequest {
	ar.request.additionalData = additionalData
	return ar
}

func (ar *AuthorizationRequest) WithCustomer(customer Customer) *AuthorizationRequest {
	ar.request.customer = &customer
	return ar
}
