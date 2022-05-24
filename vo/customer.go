package vo

import "time"

type Customer struct {
	//IP-адрес клиента, производящего оплату в вашем магазине
	Ip string `json:"ip"`
	//email клиента, производящего оплату в вашем магазине
	Email string `json:"email"`
	//id устройства клиента, производящего оплату в вашем магазине
	DeviceId string `json:"device_id"`
	//(необязательный) дата рождения клиента в формате ISO 8601 YYYY-MM-DD
	BirthDate time.Time `json:"birth_date"`
}
