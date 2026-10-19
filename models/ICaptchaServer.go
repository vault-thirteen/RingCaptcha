package m

type ICaptchaServer interface {
	Start() (err error)
	Stop() (err error)
	GetListenDsn() (dsn string)

	CreateCaptcha() (resp *CreateCaptchaResponse, err error)
	CheckCaptcha(req *CheckCaptchaRequest) (resp *CheckCaptchaResponse, err error)
	HasCaptcha(req *HasCaptchaRequest) (resp *HasCaptchaResponse, err error)
	GetCaptchaImage(req *GetCaptchaImageRequest) (resp *GetCaptchaImageResponse, err error)
}
