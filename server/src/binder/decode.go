package binder

import (
	"context"
	"net/http"
	"reflect"
)

var (
	binder Binder = &DefaultBinder{}
)

type IRequest interface {
	Validate() error
}

func BindAndValidate[T IRequest](r *http.Request) (req T, err error) {
	// Init req if it is a pointer
	reqType := reflect.TypeOf(req)
	if reqType.Kind() == reflect.Ptr {
		reqValue := reflect.ValueOf(req)
		if reqValue.IsNil() {
			req = reflect.New(reqType.Elem()).Interface().(T)
		}
		err = binder.Bind(req, r)
		if err != nil {
			return
		}
	} else {
		err = binder.Bind(&req, r)
		if err != nil {
			return
		}
	}

	err = req.Validate()
	if err != nil {
		return
	}

	return
}

// Interface for lib http decode func
func DecodeHTTPReq[T IRequest](ctx context.Context, r *http.Request) (interface{}, error) {
	return BindAndValidate[T](r)
}
