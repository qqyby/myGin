// github/pkg/errors wrapf后的error，在直接用于zap.Error(err)时候
// 会增加errorVerbose，打印栈信息
// 详见：https://github.com/uber-go/zap/issues/650
// 这里封装一下

package app

import "go.uber.org/zap"

type plainError struct {
	e error
}

func (pe plainError) Error() string {
	return pe.e.Error()
}

func PlainError(err error) zap.Field {
	return zap.Error(plainError{err})
}
