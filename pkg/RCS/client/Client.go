package client

import (
	"fmt"

	jc "github.com/ybbus/jsonrpc/v3"
)

// List of supported functions.
const (
	FuncPing          = "Ping"
	FuncCreateCaptcha = "CreateCaptcha"
	FuncCheckCaptcha  = "CheckCaptcha"
)

func NewClient(host string, port uint16, path string) jc.RPCClient {
	dsn := fmt.Sprintf("http://%s:%d%s", host, port, path)
	return jc.NewClient(dsn)
}
