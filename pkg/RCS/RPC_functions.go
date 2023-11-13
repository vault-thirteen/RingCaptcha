package rcs

import (
	"encoding/base64"
	"fmt"

	"github.com/osamingo/jsonrpc/v2"
	"github.com/vault-thirteen/RingCaptcha/pkg/RCS/models"
)

// RPC functions.

func (srv *Server) createCaptcha() (result *models.CreateCaptchaResult, jerr *jsonrpc.Error) {
	ccResponse, err := srv.captchaManager.CreateCaptcha()
	if err != nil {
		return nil, &jsonrpc.Error{Code: RpcErrorCode_CreateError, Message: fmt.Sprintf(RpcErrorMsgF_CreateError, err.Error())}
	}

	result = &models.CreateCaptchaResult{
		TaskId:              ccResponse.TaskId,
		ImageFormat:         ccResponse.ImageFormat,
		IsImageDataReturned: ccResponse.IsImageDataReturned,
	}

	if ccResponse.IsImageDataReturned {
		result.ImageDataB64 = base64.StdEncoding.EncodeToString(ccResponse.ImageData)
	}

	return result, nil
}

func (srv *Server) checkCaptcha(p *models.CheckCaptchaParams) (result *models.CheckCaptchaResult, jerr *jsonrpc.Error) {
	// Check parameters.
	if len(p.TaskId) == 0 {
		return nil, &jsonrpc.Error{Code: RpcErrorCode_TaskIdIsNotSet, Message: RpcErrorMsg_TaskIdIsNotSet}
	}

	if p.Value == 0 {
		return nil, &jsonrpc.Error{Code: RpcErrorCode_AnswerIsNotSet, Message: RpcErrorMsg_AnswerIsNotSet}
	}

	resp, err := srv.captchaManager.CheckCaptcha(&models.CheckCaptchaRequest{TaskId: p.TaskId, Value: p.Value})
	if err != nil {
		return nil, &jsonrpc.Error{Code: RpcErrorCode_CheckError, Message: fmt.Sprintf(RpcErrorMsgF_CheckError, err.Error())}
	}

	result = &models.CheckCaptchaResult{
		TaskId:    p.TaskId,
		IsSuccess: resp.IsSuccess,
	}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *models.ShowDiagnosticDataResult, jerr *jsonrpc.Error) {
	result = &models.ShowDiagnosticDataResult{
		TotalRequestsCount:      srv.diag.getTotalRequestsCount(),
		SuccessfulRequestsCount: srv.diag.getSuccessfulRequestsCount(),
	}

	return result, nil
}
