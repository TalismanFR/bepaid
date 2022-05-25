package responses

import (
	"bepaid/vo"
	"encoding/json"
	"fmt"
	"time"
)

type Customer = vo.Customer

type CreditCard struct {
	Brand         string `json:"brand"`
	Product       string `json:"product"`
	SubBrand      string `json:"sub_brand"`
	Stamp         string `json:"stamp"`
	Last4         string `json:"last_4"`
	First1        string `json:"first_1"`
	Token         string `json:"token"`
	TokenProvider string `json:"token_provider"`

	Holder   string `json:"holder"`
	ExpYear  int    `json:"exp_year"`
	ExpMonth int    `json:"exp_month"`
}

// json: "additional_data"
//секция, содержащая дополнительную информацию о платеже
type AdditionalData struct {
	//секция для работы с сервисом Masterpass
	Masterpass struct {

		//секция для параметров Masterpass
		Params struct {

			//id пользовательской сессии
			Session string `json:"session"`
		} `json:"params"`

		//результат операции в Masterpass
		Result struct {

			//статус ответа: succesfull, failed
			Status string `json:"status"`

			//сообщение о результате операции в Masterpass, сгенерированное bePaid
			Message string `json:"message"`

			//сообщение о причине ошибки в Masterpass, сгенерированное bePaid. Возвращается в случае ошибки
			Error string `json:"error"`

			//сообщение об ошибке, сгенерированное Masterpass. Возвращается в случае ошибки
			ErrorMessage string `json:"error_message"`

			//код ошибки, сгенерированный Masterpass. Возвращается в случае ошибки
			ErrorCode string `json:"error_code"`

			//token 	токен карты в системе Masterpass. Возвращается в случае сохранения карты
			Token string `json:"token"`
		} `json:"result"`
	} `json:"masterpass"`

	//секция, содержащая детальную информацию о бренде кредитной карты
	SubBrand struct {

		//название кобренда. Допустимые значения: halva
		Brand string `json:"brand"`

		//true - использовать бонусные баллы, false - не использовать бонусные баллы.
		//Параметр характерен для brand со значение halva
		UsePoints bool `json:"use_points"`
	} `json:"sub_brand"`

	//текст, который будет добавлен в письмо клиенту.
	ReceiptText string `json:"receipt_text"`
}

type BillingAddress = vo.BillingAddress

type Authorization struct {
	AuthCode          string `json:"auth_code"`
	BankCode          string `json:"bank_code"`
	Rrn               string `json:"rrn"`
	RefId             string `json:"ref_id"`
	Message           string `json:"message"`
	GatewayId         int    `json:"gateway_id"`
	BillingDescriptor string `json:"billing_descriptor"`
	Status            string `json:"status"`
}

// json: be_protected_verification
type BeProtectedVerification struct {
	Status         string `json:"status"`
	WhiteBlackList struct {
		CardNumber string `json:"card_number"`
		Ip         string `json:"ip"`
		Email      string `json:"email"`
	} `json:"white_black_list"`
	Rules map[string]any `json:"rules"`
}

type AvsCvcVerification struct {
	CvcVerification struct {
		ResultCode string `json:"result_code"`
	} `json:"cvc_verification"`
	AvsVerification struct {
		ResultCode string `json:"result_code"`
	} `json:"avs_verification"`
}

type AuthResponse struct {
	Transaction struct {
		Uid               string    `json:"uid"`
		Status            string    `json:"status"`
		Amount            int       `json:"amount"`
		Currency          string    `json:"currency"`
		Description       string    `json:"description"`
		Type              string    `json:"type"`
		TrackingId        string    `json:"tracking_id"`
		Message           string    `json:"message"`
		Test              bool      `json:"test"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		Language          string    `json:"language"`
		PaymentMethodType string    `json:"payment_method_type"`

		//ссылка на квитанцию обработанной транзакции
		ReceiptUrl string `json:"receipt_url"`

		Customer Customer `json:"customer"`

		BillingAddress BillingAddress `json:"billing_address"`

		CreditCard CreditCard `json:"credit_card"`

		Authorization Authorization `json:"authorization"`

		//необязательный блок с результатом проверки beProtected
		BeProtectedVerification map[string]any `json:"be_protected_verification"`

		//необязательный блок с результатом проверки AVS/CVC
		AvsCvcVerification map[string]any `json:"avs_cvc_verification"`

		// секция, содержащая дополнительную информацию о платеже
		AdditionalData AdditionalData `json:"additional_data"`
	} `json:"transaction"`
}

func Parse(resp string) *AuthResponse {
	var auth AuthResponse
	err := json.Unmarshal([]byte(resp), &auth)
	fmt.Println(err)
	fmt.Println(auth)

	return nil
}

var TestAuthResp1 = `
{
  "transaction":{
    "customer":{
      "ip":"127.0.0.1",
      "email":"john@example.com"
    },
    "credit_card":{
      "holder":"Johnahan Doe",
      "stamp":"a825df7faba8804619aef7a6d5a5821ec292fce04e3e43933ca33d0692df90b4",
      "brand":"visa",
      "last_4":"0000",
      "product":"Gold",
      "sub_brand":"halva",
      "first_1":"4",
      "token":"2bbd9fb7307dace37a9c2db1b4cca6f0c0dd143eac17294daab769224bff6ae2",
      "exp_month":1,
      "exp_year":2020
    },
    "receipt_url": "https://merchant.bepaid.by/customer/transactions/2-52671c8733/11443f39ae75aa1f955a9c9283cd5045bfb0413b65d666f834a9da4e7d3926b5",
      "additional_data":{
        "sub_brand":{
          "brand": "halva",
          "use_points": false
        }
      },
    "billing_address":{
      "first_name":"John",
      "last_name":"Doe",
      "address":"1st Street",
      "country":"US",
      "city":"Denver",
      "zip":"96002",
      "state":"CO",
      "phone":null
    },
    "be_protected_verification":{
      "status":"successful",
      "white_black_list":{
        "card_number":"absent",
        "ip":"absent",
        "email":"absent"
      },
      "rules":{
        "1_123_My Shop":{
          "more_100_eur" : {"Transaction amount more than 100 AND Transaction currency is EUR": "passed"}
        },
        "1_John Doe":{},
        "bePaid":{}
      }
    },
    "authorization":{
      "auth_code":"654321",
      "bank_code":"00",
      "rrn":"999",
      "ref_id":"777888",
      "message":"The operation was successfully processed.",
      "gateway_id":85,
      "billing_descriptor":"TEST GATEWAY BILLING DESCRIPTOR",
      "status":"successful"
    },
    "uid":"2-52671c8733",
    "status":"successful",
    "amount":90,
    "currency":"USD",
    "description":"Test transaction",
    "type":"authorization",
    "tracking_id":"tracking_id_000",
    "message":"Successfully processed",
    "test":true,
    "created_at":"2014-06-11T12:04:59+03:00",
    "updated_at":"2014-06-11T12:04:59+03:00",
    "avs_cvc_verification": {
      "cvc_verification" : {
        "result_code": "M"
      },
      "avs_verification" : {
        "result_code": "M"
      }
    }
  }
}
`
