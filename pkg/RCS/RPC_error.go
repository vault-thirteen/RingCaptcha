package rcs

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_TaskIdIsNotSet = 1
	RpcErrorCode_AnswerIsNotSet = 2
	RpcErrorCode_CheckError     = 3
	RpcErrorCode_CreateError    = 4
)

// Messages.
const (
	RpcErrorMsg_TaskIdIsNotSet = "task ID is not set"
	RpcErrorMsg_AnswerIsNotSet = "answer is not set"
	RpcErrorMsgF_CheckError    = "check error: %s"
	RpcErrorMsgF_CreateError   = "captcha creation error: %s"
)
