package logic

import (
	"github.com/KyoUK4n/etherscan/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

func OutSuccess(data interface{}, msg ...string) (*types.Response, error) {
	message := "success"
	if len(msg) > 0 {
		message = msg[0]
	}
	return &types.Response{
		Code: 1,
		Msg:  message,
		Data: data,
	}, nil
}

func OutFailedWithCode(errorCode int, msg ...string) (*types.Response, error) {
	message := "fail"
	if len(msg) > 0 {
		message = msg[0]
	}
	logx.Error(message)
	return &types.Response{
		Code: errorCode,
		Msg:  message,
		Data: nil,
	}, nil
}

func OutFailed(msg ...string) (*types.Response, error) {
	message := "fail"
	if len(msg) > 0 {
		message = msg[0]
	}
	return &types.Response{
		Code: 0,
		Msg:  message,
		Data: nil,
	}, nil
}

func OutFailedWithErr(err error, msg ...string) (*types.Response, error) {
	message := "fail"
	if len(msg) > 0 {
		message = msg[0]
	}
	logx.Errorf("%s, cause by: %v", message, err)
	return &types.Response{
		Code: 0,
		Msg:  message,
		Data: nil,
	}, nil
}
