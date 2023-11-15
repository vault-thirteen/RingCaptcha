package rc

type CreateCaptchaResponse struct {
	TaskId              string
	ImageFormat         string
	IsImageDataReturned bool
	ImageData           []byte
}

type CheckCaptchaRequest struct {
	TaskId string
	Value  uint
}

type CheckCaptchaResponse struct {
	TaskId    string
	IsSuccess bool
}
