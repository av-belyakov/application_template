package wrappers

import (
	"fmt"
	"runtime"
)

// WrapperError ошибка оборачивается ссылкой на имя файла с кодом и номером строки в коде
func WrapperError(err error) error {
	_, f, l, _ := runtime.Caller(1)

	return fmt.Errorf("%w %s:%d", err, f, l)
}
