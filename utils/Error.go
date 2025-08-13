package utils

import (
	"fmt"
	"runtime"

	"github.com/mytionbr/temperature_system/model"
)

func NewError(err error, code int) model.StatusError {
	pc, _, line, _ := runtime.Caller(2)
	details := runtime.FuncForPC(pc)

	return model.StatusError{
		Code:    code,
		Err:     err,
		Message: fmt.Sprintf("%v", err),
		Caller:  fmt.Sprintf("%s#%d", details.Name(), line),
	}
}
