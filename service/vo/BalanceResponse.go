package vo

type BalanceResponse struct {
	Status string `json:"status"`
	Result struct {
		GatewayId int    `json:"gatewayId"`
		Account   string `json:"account"`
		Amount    int64  `json:"amount"`
		Currency  string `json:"currency"`
		BankInfo  struct {
			Account string  `json:"Account"`
			Amount  float64 `json:"Amount"`
			Balance struct {
				OperDate     string  `json:"OperDate"`
				Credit       float64 `json:"Credit"`
				CreditRub    float64 `json:"CreditRub"`
				Debit        float64 `json:"Debit"`
				DebitRub     float64 `json:"DebitRub"`
				AmountIn     float64 `json:"AmountIn"`
				AmountInRub  float64 `json:"AmountInRub"`
				AmountOut    float64 `json:"AmountOut"`
				AmountOutRub float64 `json:"AmountOutRub"`
			} `json:"Balance"`
		} `json:"bankInfo"`
	} `json:"result"`
}

func (tr *BalanceResponse) IsSuccess() bool {
	return tr.Status == success
}
