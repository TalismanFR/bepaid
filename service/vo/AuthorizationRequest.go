package vo

import (
	"encoding/json"
	"fmt"
)

type AuthorizationRequest struct {

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

func (ar AuthorizationRequest) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		Request interface{} `json:"request"`
	}{struct {
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
	}{
		ar.amount,
		ar.currency,
		ar.description,
		ar.trackingId,
		ar.duplicateCheck,
		ar.returnUrl,
		ar.test,
		ar.creditCard,
		ar.additionalData,
		ar.customer,
	}})
}

func (ar *AuthorizationRequest) Amount() Amount {
	return ar.amount
}

func (ar *AuthorizationRequest) Currency() Currency {
	return ar.currency
}

func (ar *AuthorizationRequest) Description() string {
	return ar.description
}

func (ar *AuthorizationRequest) TrackingId() string {
	return ar.trackingId
}

func (ar *AuthorizationRequest) DuplicateCheck() *bool {
	return ar.duplicateCheck
}

func (ar *AuthorizationRequest) ReturnUrl() string {
	return ar.returnUrl
}

func (ar *AuthorizationRequest) Test() bool {
	return ar.test
}

func (ar *AuthorizationRequest) CreditCard() CreditCard {
	return ar.creditCard
}

func (ar *AuthorizationRequest) AdditionalData() map[string]interface{} {
	return ar.additionalData
}

func (ar *AuthorizationRequest) Customer() *Customer {
	return ar.customer
}

// NewAuthorizationRequest creates AuthorizationRequest with mandatory fields
func NewAuthorizationRequest(amount Amount, currency Currency, description string, trackingId string, test bool, cc CreditCard) *AuthorizationRequest {

	r := &AuthorizationRequest{}

	r.amount = amount
	r.currency = currency
	r.description = description
	r.trackingId = trackingId
	r.test = test
	r.creditCard = cc

	return r
}

func (ar *AuthorizationRequest) Valid() (bool, error) {
	if len(ar.description) > 255 {
		return false, fmt.Errorf("description length greater than 255")
	}
	if len(ar.trackingId) > 255 {
		return false, fmt.Errorf("trackingId length greater than 255")
	}

	return true, nil
}

func (ar *AuthorizationRequest) WithDuplicateCheck(duplicateCheck bool) *AuthorizationRequest {
	ar.duplicateCheck = &duplicateCheck
	return ar
}

func (ar *AuthorizationRequest) WithReturnUrl(returnUrl string) *AuthorizationRequest {
	ar.returnUrl = returnUrl
	return ar
}

// WithAdditionalData saves argument to Authorizationrequest.request.AdditionalData field.
//
// Don't change content of additionalData after function call.
func (ar *AuthorizationRequest) WithAdditionalData(additionalData map[string]interface{}) *AuthorizationRequest {
	ar.additionalData = additionalData
	return ar
}

func (ar *AuthorizationRequest) WithCustomer(customer Customer) *AuthorizationRequest {
	ar.customer = &customer
	return ar
}
