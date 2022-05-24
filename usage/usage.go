package main

import (
	"bepaid/response"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	response.Parse(response.TestAuthResp1)

	//baseURL := "https://gateway.bepaid.by/transactions/"
	//shopId := ""
	//secret := ""
	//
	//httpClient := &http.Client{Transport: &mockTransport{},}
	////httpClient := http.DefaultClient
	//
	//api1 := bepaid.NewApi(httpClient, baseURL, shopId, secret)
	//
	//client1 := bepaid.NewClient(api1)
	//
	//_, err := client1.Capture(context.Background(), vo.CaptureRequest{ParentUid: "test1", Amount: 345, DuplicateCheck: true})
	//if err != nil {
	//	fmt.Println(err)
	//}

	//_, err := client1.Authorizations(context.Background(), vo.AuthorizationRequest{})
	//if err != nil {
	//	fmt.Println(err)
	//}

}

type mockTransport struct {
}

func (m *mockTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	fmt.Println("here")
	ma := map[string]any{}
	if err := json.NewDecoder(request.Body).Decode(&ma); err != nil {
		log.Fatal(err)
	}

	fmt.Println(ma)

	return nil, nil
}
