package vo

type TokenizationRequest struct {
	Request struct {

		//стоимость в минимальных денежных единицах.
		//Например, $32.45 должна быть отправлена как 3245
		Amount int64 `json:"amount"`

		//валюта в ISO-4217 формате, например USD
		Currency string `json:"currency"`

		//описание заказа. Максимальная длина: 255 символов
		Description string `json:"description"`

		//id транзакции или заказа в вашей системе.
		//Максимальная длина: 255 символов.
		//Пожалуйста, используйте уникальное значение для того, чтобы при запросе статуса транзакции получить актуальную информацию.
		//В противном случае вы получите первую найденную по tracking_id транзакцию
		TrackingId string `json:"tracking_id"`

		//(необязательный) язык вашей страницы оформления заказа. Если
		//параметр установлен и email уведомление о транзакции включено,
		//то bePaid отправит email, язык текста которого будет language. По
		//умолчанию - en. Допустимые значения:
		Language string `json:"language,omitempty"`

		//(необязательный) узнайте у службы поддержки, должны ли вы отправлять эти данные
		BillingAddress BillingAddress `json:"billing_address,omitempty"`

		CreditCard CreditCard `json:"credit_card"`
		Customer   Customer   `json:"customer"`

		//секция, содержащая дополнительную информацию о платеже
		AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
	} `json:"request"`
}

func NewTokenizationRequest(amount int64, currency, description, trackingId string, creditCard CreditCard, customer Customer) *TokenizationRequest {
	r := &TokenizationRequest{}

	r.Request.Amount = amount
	r.Request.Currency = currency
	r.Request.Description = description
	r.Request.TrackingId = trackingId
	r.Request.CreditCard = creditCard
	r.Request.Customer = customer

	return r
}

func (tr *TokenizationRequest) WithAdditionalData(additionalData map[string]interface{}) *TokenizationRequest {
	tr.Request.AdditionalData = additionalData
	return tr
}

func (tr *TokenizationRequest) WithBillingAddress(billingAddress BillingAddress) *TokenizationRequest {
	tr.Request.BillingAddress = billingAddress
	return tr
}

func (tr *TokenizationRequest) WithLanguage(language string) *TokenizationRequest {
	tr.Request.Language = language
	return tr
}
