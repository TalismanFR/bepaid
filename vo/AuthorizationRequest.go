package vo

type AuthorizationRequest struct {
}

func (a AuthorizationRequest) Read(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func NewAuthorizationRequest() *AuthorizationRequest {
	return &AuthorizationRequest{}
}
