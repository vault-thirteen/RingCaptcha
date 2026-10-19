package es

import (
	"github.com/vault-thirteen/RingCaptcha/models"
	"github.com/vault-thirteen/RingCaptcha/server"
)

type EasyServer struct {
	css *m.CaptchaServerSettings
	cs  *s.CaptchaServer
}

func NewEasyServer(css *m.CaptchaServerSettings) (es *EasyServer) {
	es = &EasyServer{
		css: css,
	}

	var err error
	es.cs, err = s.NewCaptchaServer(css)
	mustBeNoError(err)
	return es
}

func (es *EasyServer) Start() {
	var err error
	err = es.cs.Start()
	mustBeNoError(err)
}
func (es *EasyServer) Stop() {
	var err error
	err = es.cs.Stop()
	mustBeNoError(err)
}
func (es *EasyServer) GetListenDsn() (dsn string) {
	return es.cs.GetListenDsn()
}

func (es *EasyServer) CreateCaptcha() (captcha *m.Captcha) {
	resp, err := es.cs.CreateCaptcha()
	mustBeNoError(err)

	captcha = &m.Captcha{
		Id:          resp.TaskId,
		ImageData:   resp.ImageData,
		ImageFormat: resp.ImageFormat,
	}

	return captcha
}
func (es *EasyServer) CheckCaptcha(captcha *m.Captcha) (ok bool) {
	var resp *m.CheckCaptchaResponse
	var err error
	resp, err = es.cs.CheckCaptcha(&m.CheckCaptchaRequest{TaskId: captcha.Id, Value: captcha.RingCount})
	mustBeNoError(err)

	return resp.IsSuccess
}
func (es *EasyServer) HasCaptcha(captcha *m.Captcha) (ok bool) {
	var resp *m.HasCaptchaResponse
	var err error
	resp, err = es.cs.HasCaptcha(&m.HasCaptchaRequest{TaskId: captcha.Id})
	mustBeNoError(err)

	return resp.IsFound
}
func (es *EasyServer) GetCaptchaImage(captcha *m.Captcha) (data []byte) {
	var resp *m.GetCaptchaImageResponse
	var err error
	resp, err = es.cs.GetCaptchaImage(&m.GetCaptchaImageRequest{TaskId: captcha.Id})
	mustBeNoError(err)

	return resp.ImageData
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
