package common

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonerror "github.com/chenjun-git/umbrella-common/proto"
)

// Error 兼容用的错误
type Error struct {
	commonError *commonerror.Error
	error       error
}

// Errorf 生成 *common.Error
func Errorf(code int, format string, a ...interface{}) *Error {
	// todo 1.本修改上线2.修改tenant-gateway3.直接返回 *commonerror.Error
	return &Error{
		commonError: &commonerror.Error{
			Code:        int32(code),
			Description: fmt.Sprintf(format, a...),
		},
		error: status.Errorf(codes.Code(code), fmt.Sprintf("[%d: %s] %s", code, ErrorMap[code]["ZH-CN"], format), a...),
	}
}

// GetError 获取 error 类型的
func (e *Error) GetError() error {
	if e == nil {
		return nil
	}

	return e.error
}

// GetCommonError 获取 *commonerror.Error 类型的
func (e *Error) GetCommonError() *commonerror.Error {
	if e == nil {
		return nil
	}

	return e.commonError
}

// NormalErrorStatus 返回错误码对应的 error
func NormalErrorStatus(code int) error {
	return status.Error(codes.Code(code), GetMsg(code, []string{"ZH-CN"}))
}

// ExtendErrorStatus 返回错误码+自定义错误信息对应的 error
func ExtendErrorStatus(code int, err error) error {
	if err == nil {
		err = errors.New("")
	}
	return status.Errorf(codes.Code(code), "%s:%s", GetMsg(code, []string{"ZH-CN"}), err.Error())
}