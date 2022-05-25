package responses

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Example of responses
//"responses": {
//        "message": "Amount is not a number. Amount is not a number. Description can't be blank. Currency can't be blank. Exp year is invalid (format: yyyy). Exp month is invalid (format: mm, only 01-12 allowed). Number is invalid. Verification value Only digits allowed,Card verification value (CVC/CVV) is invalid. It should be 3 digits. Date is expired,can't be blank.",
//        "errors": {
//            "amount": [
//                "is not a number",
//                "is not a number"
//            ],
//            "description": [
//                "can't be blank"
//            ],
//            "currency": [
//                "can't be blank"
//            ],
//            "credit_card": {
//                "exp_year": [
//                    "is invalid (format: yyyy)"
//                ],
//                "exp_month": [
//                    "is invalid (format: mm, only 01-12 allowed)"
//                ],
//                "number": [
//                    "is invalid"
//                ],
//                "verification_value": [
//                    "Only digits allowed",
//                    "Card verification value (CVC/CVV) is invalid. It should be 3 digits"
//                ],
//                "date": [
//                    "is expired",
//                    "can't be blank"
//                ]
//            }
//        }
//    }

// Response contains info about errors in request
type Response struct {

	// Message is a string containing all error messages from Errors
	Message string `json:"message"`

	// Errors is a map containing errors specific for every field of every object
	//
	// Errors should be checked for nil
	Errors map[string]any `json:"errors"`
}

func (r *Response) Error() string {
	buf := bytes.Buffer{}
	fmt.Fprintln(&buf, r.Message)
	b, _ := json.MarshalIndent(r.Errors, "", "\t")
	buf.Write(b)
	return string(buf.Bytes())
}

func Is(target error) bool {
	return true
}

//struct {
//		Amount      []string `json:"amount"`
//		Description []string `json:"description"`
//		Currency    []string `json:"currency"`
//		CreditCard  struct {
//			ExpYear           []string `json:"exp_year"`
//			ExpMonth          []string `json:"exp_month"`
//			Number            []string `json:"number"`
//			VerificationValue []string `json:"verification_value"`
//			Date              []string `json:"date"`
//		} `json:"credit_card"`
//	}
