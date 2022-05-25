package vo

type AuthorizationRequest struct {

	//стоимость в минимальных денежных единицах.
	//Например, $32.45 должна быть отправлена как 3245
	Amount int `json:"amount"`

	//валюта в ISO-4217 формате, например USD
	Currency string `json:"currency"`

	//описание заказа. Максимальная длина: 255 символов
	Description string `json:"description,omitempty"`

	//id транзакции или заказа в вашей системе.
	//Максимальная длина: 255 символов.
	//Пожалуйста, используйте уникальное значение для того, чтобы при запросе статуса транзакции получить актуальную информацию.
	//В противном случае вы получите первую найденную по tracking_id транзакцию
	TrackingId string `json:"tracking_id"`

	//параметр обязателен, если 3-D Secure включен.
	//Обратитесь к менеджеру за информацией. return_url - это URL на стороне торговца,
	//на который bePaid будет перенаправлять клиента после возврата с 3-D Secure проверки
	ReturnUrl string `json:"return_url,omitempty"`

	//true или false. Транзакция будет тестовой, если значение true.
	Test bool `json:"test"`

	CreditCard CreditCard `json:"credit_card"`

	//секция, содержащая дополнительную информацию о платеже
	AdditionalData map[string]any `json:"additional_data,omitempty"`

	Customer *Customer `json:"customer,omitempty"`

	//(необязательный) узнайте у службы поддержки, должны ли вы отправлять эти данные
	BillingAddress *BillingAddress `json:"billing_address,omitempty"`
}

func NewAuthorizationRequest() *AuthorizationRequest {
	return &AuthorizationRequest{}
}

// NewAuthorizationRequestWith creates AuthorizationRequest with mandatory fields
func NewAuthorizationRequestWith(amount int, currency string, description string, cc CreditCard) *AuthorizationRequest {
	return &AuthorizationRequest{}
}

func (a *AuthorizationRequest) WithTrackingId(trackingId string) *AuthorizationRequest {
	a.TrackingId = trackingId
	return a
}

func (a *AuthorizationRequest) WithReturnUrl(returnUrl string) *AuthorizationRequest {
	a.ReturnUrl = returnUrl
	return a
}

func (a *AuthorizationRequest) WithTest(test bool) *AuthorizationRequest {
	a.Test = test
	return a
}

func (a *AuthorizationRequest) WithAdditionalData(additionalData map[string]any) *AuthorizationRequest {
	a.AdditionalData = additionalData
	return a
}

func (a *AuthorizationRequest) WithCustomer(customer Customer) *AuthorizationRequest {
	*a.Customer = customer
	return a
}

func (a *AuthorizationRequest) WithBillingAddress(billingAddress BillingAddress) *AuthorizationRequest {
	*a.BillingAddress = billingAddress
	return a
}
