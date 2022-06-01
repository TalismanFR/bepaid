package vo

type BillingAddress struct {
	//имя клиента. Максимальная длина: 30 символов
	FirstName string `json:"first_name"`
	//фамилия клиента. Максимальная: 30 символов
	LastName string `json:"last_name"`
	//страна клиента в ISO 3166-1 alpha-2 формате
	Country string `json:"country"`
	//город клиента. Максимальная длина: 60 символов
	City string `json:"city"`
	//двухбуквенная абревиатура штата, если страна клиента US или CA
	State string `json:"state"`
	//(необязательный) почтовый индекс клиента. Для country=US, формат
	//почтового индекса должен иметь вид NNNNN или NNNNN-NNNN
	Zip string `json:"zip,omitempty"`
	//адрес клиента. Максимальная длина: 255 символов
	Address string `json:"address"`
}

func NewBillingAddress(firstName, lastName, country, city, state, zip, address string) *BillingAddress {
	return &BillingAddress{FirstName: firstName, LastName: lastName, Country: country, City: city, State: state, Zip: zip, Address: address}
}
