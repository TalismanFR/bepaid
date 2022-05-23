package vo

type CaptureRequest struct {
}

func (c CaptureRequest) Read(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func NewCaptureRequest() *CaptureRequest {
	return &CaptureRequest{}
}
