package usage

import (
	"bepaid"
	"bepaid/vo"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	baseURL := "https://gateway.bepaid.by/transactions/"
	shopId := ""
	secret := ""

	httpClient := http.DefaultClient

	api1 := bepaid.NewApi(httpClient, baseURL, shopId, secret)

	client1 := bepaid.NewClient(api1)

	r := vo.AuthorizationRequest{}
	err := json.Unmarshal([]byte(Test1), &r)
	if err != nil {
		log.Fatal(err)
	}

	tr, err := client1.Authorizations(context.Background(), r)
	if err != nil {
		fmt.Println("Response")
		fmt.Println(err)
	}

	fmt.Println("Transaction")
	fmt.Println(tr)

}

var Test1 = `
{
      "amount":100,
      "currency":"USD",
      "description":"Test transaction",
      "tracking_id":"your_uniq_number",
      "language":"en",
      "test":true,
      "billing_address":{
         "first_name":"John",
         "last_name":"Doe",
         "country":"US",
         "city":"Denver",
         "state":"CO",
         "zip":"96002",
         "address":"1st Street"
      },
      "credit_card":{
         "number":"4200000000000000",
         "verification_value":"123",
         "holder":"John Doe",
         "exp_month":"05",
         "exp_year":"2020"
      },
      "additional_data":{
         "sub_brand":{
            "brand": "halva",
            "use_points": false
         }
      },
      "customer":{
         "ip":"127.0.0.1",
         "email":"john@example.com"
      }
   }
`
