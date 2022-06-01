package vo

type PayoutRequest struct {
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

		//(необязательный) узнайте у службы поддержки, должны ли вы отправлять эти данные
		RecipientBillingAddress BillingAddress `json:"recipient_billing_address,omitempty"`

		//(необязательный) узнайте у службы поддержки, должны ли вы отправлять эти данные
		SenderBillingAddress BillingAddress `json:"sender_billing_address,omitempty"`

		RecipientCreditCard CreditCard `json:"recipient_credit_card"`

		//секция, содержащая дополнительную информацию о платеже
		AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
	} `json:"request"`
}

func NewPayoutRequest(amount int64, currency, description, trackingId string, creditCard CreditCard) *TokenizationRequest {
	r := &TokenizationRequest{}

	r.Request.Amount = amount
	r.Request.Currency = currency
	r.Request.Description = description
	r.Request.TrackingId = trackingId
	r.Request.CreditCard = creditCard

	return r
}

func (pr *PayoutRequest) WithRecipientBillingAddress(recipientBillingAddress BillingAddress) *PayoutRequest {
	pr.Request.RecipientBillingAddress = recipientBillingAddress
	return pr
}

func (pr *PayoutRequest) WithSenderBillingAddress(senderBillingAddress BillingAddress) *PayoutRequest {
	pr.Request.SenderBillingAddress = senderBillingAddress
	return pr
}

func (pr *PayoutRequest) WithAdditionalData(additionalData map[string]interface{}) *PayoutRequest {
	pr.Request.AdditionalData = additionalData
	return pr
}
