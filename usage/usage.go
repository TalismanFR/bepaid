package usage

import (
	"bepaid"
	"net/http"
)

func usage() {
	baseURL := "https://gateway.bepaid.by/transactions/"
	shop_id := ""
	secret := ""

	api1 := bepaid.NewApi(http.DefaultClient, baseURL, shop_id, secret)

	client1 := bepaid.NewClient(api1)

	_ = client1

}
