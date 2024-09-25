package errs

import (
	"errors"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

func WrapMessage(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func ExtractStatusCode(err error) (code int) {
	if err == nil {
		return http.StatusBadRequest
	}

	code = http.StatusBadRequest
	for {
		if err == nil {
			break
		}
		if e, ok := err.(kithttp.StatusCoder); ok {
			code = e.StatusCode()
			break
		}
		err = errors.Unwrap(err)
	}
	return code
}
