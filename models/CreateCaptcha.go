package m

type CreateCaptchaRequest struct{}

type CreateCaptchaResponse struct {
	TaskId      string
	ImageFormat string

	// When storage is used, images are not returned in response objects;
	// instead, they are manually requested by clients from storage. Opposite
	// is also true: when storage is not used, server returns images in
	// response objects.
	IsImageDataReturned bool
	ImageData           []byte
}
