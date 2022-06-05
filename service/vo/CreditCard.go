package vo

import (
	"encoding/json"
	"fmt"
)

// CreditCard
//
// Use NewCreditCardWithToken if you already have card token
type CreditCard struct {

	// номер карты, длина - от 12 до 19 цифр
	number string `json:"number"`

	//3-х или 4-х цифровой код безопасности (CVC2, CVV2 или CID, в зависимости от бренда карты).
	//Может быть отправлен вместе с параметром token и bePaid доставит банку-эквайеру данные карты с CVC2/CVV2/CID
	verificationValue string `json:"verification_value"`

	//имя владельца карты. Максимальная длина: 32 символа
	holder string `json:"holder"`

	//месяц окончания срока действия карты, представленный двумя цифрами (например, 01)
	expMonth ExpMonth `json:"exp_month"`

	//год срока окончания действия карты, представленный четырьмя цифрами (например, 2007)
	expYear string `json:"exp_year"`

	//опционально
	//вместо 5 параметров выше можно отправить токен карты, который был получен в ответе первой оплаты.
	//Если используется токен карты, то необходимо обязательно указывать параметр additional_data.contract
	token string `json:"token,omitempty"`

	//опционально
	//если значение параметра true, bePaid не выполняет 3-D Secure проверку.
	//Это полезно если вы, например, не хотите чтобы клиент проходил 3-D Secure проверку снова.
	//Уточните у службы поддержки, можете ли вы использовать этот параметр.
	skipThreeDSecureVerification bool `json:"skip_three_d_secure_verification"`
}

func (cc CreditCard) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Number                       string   `json:"number"`
		VerificationValue            string   `json:"verification_value"`
		Holder                       string   `json:"holder"`
		ExpMonth                     ExpMonth `json:"exp_month"`
		ExpYear                      string   `json:"exp_year"`
		Token                        string   `json:"token,omitempty"`
		SkipThreeDSecureVerification bool     `json:"skip_three_d_secure_verification"`
	}{
		cc.number,
		cc.verificationValue,
		cc.holder,
		cc.expMonth,
		cc.expYear,
		cc.token,
		cc.skipThreeDSecureVerification,
	})
}

func (cc *CreditCard) Number() string {
	return cc.number
}

func (cc *CreditCard) VerificationValue() string {
	return cc.verificationValue
}

func (cc *CreditCard) Holder() string {
	return cc.holder
}

func (cc *CreditCard) ExpMonth() ExpMonth {
	return cc.expMonth
}

func (cc *CreditCard) ExpYear() string {
	return cc.expYear
}

func (cc *CreditCard) Token() string {
	return cc.token
}

func (cc *CreditCard) SkipThreeDSecureVerification() bool {
	return cc.skipThreeDSecureVerification
}

func NewCreditCard(number string, verificationValue string, holder string, expMonth ExpMonth, expYear string) *CreditCard {
	return &CreditCard{
		number:            number,
		verificationValue: verificationValue,
		holder:            holder,
		expMonth:          expMonth,
		expYear:           expYear,
	}
}

func (cc *CreditCard) Valid() (bool, error) {
	if len(cc.number) < 12 || len(cc.number) > 19 {
		return false, fmt.Errorf("number length should be between 12 and 19 (both including)")
	}
	if len(cc.verificationValue) != 3 && len(cc.verificationValue) != 4 {
		return false, fmt.Errorf("verificationCode length should equal to 3 or 4")
	}
	if len(cc.holder) > 32 {
		return false, fmt.Errorf("holder length greater than 32")
	}
	if len(cc.expMonth) != 2 {
		return false, fmt.Errorf("expMonth length should equal 2")
	}
	if len(cc.expYear) != 4 {
		return false, fmt.Errorf("expYear length should equal 4")
	}

	return true, nil
}

func NewCreditCardWithToken(token string) *CreditCard {
	return &CreditCard{token: token}
}

func (cc *CreditCard) WithSkip3DSVerification(skipThreeDSecureVerification bool) *CreditCard {
	cc.skipThreeDSecureVerification = skipThreeDSecureVerification
	return cc
}
