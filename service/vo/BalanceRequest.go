package vo

type BalanceRequest struct {
	Request struct {
		//Номер счета (если поле пустое, то вернется информация о счете по умолчанию)
		Account string `json:"account,omitempty"`

		//Валюта в ISO-4217 формате, например USD
		Currency string `json:"currency,omitempty"`

		//ID банка, шлюза в системе bePaid
		GatewayId string `json:"gateway_id,omitempty"`
	} `json:"request"`
}

func NewBalanceRequest(account string) *BalanceRequest {
	r := &BalanceRequest{}

	r.Request.Account = account
	return r
}

func (br *BalanceRequest) WithCurrency(currency string) *BalanceRequest {
	br.Request.Currency = currency
	return br
}

func (br *BalanceRequest) WithGatewayId(gatewayId string) *BalanceRequest {
	br.Request.GatewayId = gatewayId
	return br
}
