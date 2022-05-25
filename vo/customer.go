package vo

type Customer struct {

	//IP-адрес клиента, производящего оплату в вашем магазине
	Ip string `json:"ip"`

	//email клиента, производящего оплату в вашем магазине
	Email string `json:"email"`

	//id устройства клиента, производящего оплату в вашем магазине
	//необязательно (из примеров запросов)
	DeviceId string `json:"device_id"`

	//(необязательный) дата рождения клиента в формате ISO 8601 YYYY-MM-DD
	BirthDate string `json:"birth_date"`
}

func NewCustomer(ip, email string) *Customer {
	return &Customer{Ip: ip, Email: email}
}

func (c *Customer) WithDeviceId(deviceId string) *Customer {
	c.DeviceId = deviceId
	return c
}

func (c *Customer) WithBirthDate(birthDate string) *Customer {
	c.BirthDate = birthDate
	return c
}
