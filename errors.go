package bepaid

import "fmt"

var (
	// 401
	ErrIvalidAuth = fmt.Errorf("invalid authorization: check username and password")
	// 400
	ErrBadRequest = fmt.Errorf("bad request: check responses body for more datails")
)
