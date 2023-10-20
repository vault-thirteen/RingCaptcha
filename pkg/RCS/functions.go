package rcs

import "github.com/vault-thirteen/RingCaptcha/pkg/RCS/models"

// Server functions.

func (srv *Server) createCaptcha() (resp *models.CreateCaptchaResponse, err error) {
	srv.cmGuard.Lock()
	defer srv.cmGuard.Unlock()

	resp, err = srv.captchaManager.CreateCaptcha()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (srv *Server) checkCaptcha(req *models.CheckCaptchaRequest) (resp *models.CheckCaptchaResponse, err error) {
	srv.cmGuard.Lock()
	defer srv.cmGuard.Unlock()

	resp, err = srv.captchaManager.CheckCaptcha(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
