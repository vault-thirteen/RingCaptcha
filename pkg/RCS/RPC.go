package rcs

// RPC handlers.

import (
	"context"
	"encoding/json"
	"time"

	"github.com/osamingo/jsonrpc/v2"
	"github.com/vault-thirteen/RingCaptcha/pkg/RCS/client"
	"github.com/vault-thirteen/RingCaptcha/pkg/RCS/models"
)

func (srv *Server) initJsonRpcHandlers() (err error) {
	err = srv.jsonRpcHandlers.RegisterMethod(client.FuncPing, PingHandler{Server: srv}, models.PingParams{}, models.PingResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(client.FuncCreateCaptcha, CreateCaptchaHandler{Server: srv}, models.CreateCaptchaParams{}, models.CreateCaptchaResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(client.FuncCheckCaptcha, CheckCaptchaHandler{Server: srv}, models.CheckCaptchaParams{}, models.CheckCaptchaResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(client.FuncShowDiagnosticData, ShowDiagnosticDataHandler{Server: srv}, models.ShowDiagnosticDataParams{}, models.ShowDiagnosticDataResult{})
	if err != nil {
		return err
	}

	return nil
}

type PingHandler struct {
	Server *Server
}

func (h PingHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *jsonrpc.Error) {
	h.Server.diag.incTotalRequestsCount()
	result := models.PingResult{OK: true}
	h.Server.diag.incSuccessfulRequestsCount()
	return result, nil
}

type CreateCaptchaHandler struct {
	Server *Server
}

func (h CreateCaptchaHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *jsonrpc.Error) {
	h.Server.diag.incTotalRequestsCount()
	var timeStart = time.Now()

	result, jerr := h.Server.createCaptcha()
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.incSuccessfulRequestsCount()
	return result, nil
}

type CheckCaptchaHandler struct {
	Server *Server
}

func (h CheckCaptchaHandler) ServeJSONRPC(c context.Context, params *json.RawMessage) (any, *jsonrpc.Error) {
	h.Server.diag.incTotalRequestsCount()
	var timeStart = time.Now()

	var p models.CheckCaptchaParams
	jerr := jsonrpc.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *models.CheckCaptchaResult
	result, jerr = h.Server.checkCaptcha(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.incSuccessfulRequestsCount()
	return result, nil
}

type ShowDiagnosticDataHandler struct {
	Server *Server
}

func (h ShowDiagnosticDataHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *jsonrpc.Error) {
	h.Server.diag.incTotalRequestsCount()
	var timeStart = time.Now()

	result, jerr := h.Server.showDiagnosticData()
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.incSuccessfulRequestsCount()
	return result, nil
}
