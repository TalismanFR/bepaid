package vo

import (
	"bepaid/enums"
	"fmt"
)

type CreditCard struct {

	// номер карты, длина - от 12 до 19 цифр
	Number Number `json:"number"`

	//3-х или 4-х цифровой код безопасности (CVC2, CVV2 или CID, в зависимости от бренда карты).
	//Может быть отправлен вместе с параметром token и bePaid доставит банку-эквайеру данные карты с CVC2/CVV2/CID
	VerificationValue string `json:"verification_value"`

	//имя владельца карты. Максимальная длина: 32 символа
	Holder string `json:"holder"`

	//месяц окончания срока действия карты, представленный двумя цифрами (например, 01)
	ExpMonth enums.ExpMonth `json:"exp_month"`

	//год срока окончания действия карты, представленный четырьмя цифрами (например, 2007)
	ExpYear ExpYear `json:"exp_year"`

	//опционально
	//вместо 5 параметров выше можно отправить токен карты, который был получен в ответе первой оплаты.
	//Если используется токен карты, то необходимо обязательно указывать параметр additional_data.contract
	Token string `json:"token"`

	//опционально
	//если значение параметра true, bePaid не выполняет 3-D Secure проверку.
	//Это полезно если вы, например, не хотите чтобы клиент проходил 3-D Secure проверку снова.
	//Уточните у службы поддержки, можете ли вы использовать этот параметр.
	SkipThreeDSecureVerification bool `json:"skip_three_d_secure_verification"`
}

type ExpYear string
type Number string

func NewNumber(number string) (Number, error) {
	if len(number) < 12 || len(number) > 19 {
		return "", fmt.Errorf("invalid Number: string length should be between 12 and 19")
	}
	return Number(number), nil
}
func NewExpYear(year string) (ExpYear, error) {
	if len(year) != 4 {
		return "", fmt.Errorf("invalid ExpYear: string lenght should be exactly 4")
	}
	return ExpYear(year), nil
}
