package rcs

// RPC handlers.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/osamingo/jsonrpc/v2"
	"github.com/vault-thirteen/RingCaptcha/pkg/RCS/client"
	"github.com/vault-thirteen/RingCaptcha/pkg/RCS/models"
)

const (
	RpcErrorCodeCreateCaptcha = 1
	RpcErrorCodeCheckCaptcha  = 2
)

func (srv *Server) initJsonRpcHandlers() (err error) {
	err = srv.jsonRpcHandlers.RegisterMethod(client.FuncPing, PingHandler{}, models.PingParams{}, models.PingResult{})
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

	return nil
}

type PingHandler struct{}

func (h PingHandler) ServeJSONRPC(c context.Context, params *json.RawMessage) (any, *jsonrpc.Error) {
	return models.PingResult{OK: true}, nil
}

type CreateCaptchaHandler struct {
	Server *Server
}

func (h CreateCaptchaHandler) ServeJSONRPC(c context.Context, params *json.RawMessage) (any, *jsonrpc.Error) {
	var timeStart = time.Now()

	srvResponse, err := h.Server.createCaptcha()
	if err != nil {
		return nil, &jsonrpc.Error{
			Code:    RpcErrorCodeCreateCaptcha,
			Message: err.Error(),
		}
	}

	rpcResponse := &models.CreateCaptchaResult{
		TaskId:              srvResponse.TaskId,
		ImageFormat:         srvResponse.ImageFormat,
		IsImageDataReturned: srvResponse.IsImageDataReturned,
	}

	if rpcResponse.IsImageDataReturned {
		rpcResponse.ImageDataB64 = base64.StdEncoding.EncodeToString(srvResponse.ImageData)
	}

	rpcResponse.TimeSpent = time.Now().Sub(timeStart).Milliseconds()

	return rpcResponse, nil
}

type CheckCaptchaHandler struct {
	Server *Server
}

func (h CheckCaptchaHandler) ServeJSONRPC(c context.Context, params *json.RawMessage) (any, *jsonrpc.Error) {
	var p models.CheckCaptchaParams
	jerr := jsonrpc.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var timeStart = time.Now()

	srvResponse, err := h.Server.checkCaptcha(
		&models.CheckCaptchaRequest{
			TaskId: p.TaskId,
			Value:  p.Value,
		},
	)
	if err != nil {
		return nil, &jsonrpc.Error{
			Code:    RpcErrorCodeCheckCaptcha,
			Message: err.Error(),
		}
	}

	rpcResponse := &models.CheckCaptchaResult{
		TaskId:    srvResponse.TaskId,
		IsSuccess: srvResponse.IsSuccess,
	}

	rpcResponse.TimeSpent = time.Now().Sub(timeStart).Milliseconds()

	return rpcResponse, nil
}
