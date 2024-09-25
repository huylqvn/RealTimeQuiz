package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"quizserver/src/errs"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
	"go.elastic.co/apm/v2"
)

// encodeJSONResponse is the common method to encode all response types to the client.
func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return kithttp.EncodeJSONResponse(ctx, w, response)
}

func streamMinioObject(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	obj := response.(*minio.Object)
	defer obj.Close()

	stat, err := obj.Stat()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", stat.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size))

	// Use inline to prevent download, show in browser
	// w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", stat.Key))
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", strings.ReplaceAll(stat.Key, "/", "_")))

	w.Header().Set("Last-Modified", stat.LastModified.UTC().Format(http.TimeFormat))
	for k, v := range stat.Metadata {
		w.Header().Set(k, strings.Join(v, ", "))
	}

	_, err = io.Copy(w, obj)
	return err
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	// maybe we can be smart here by returning text/json error based on request's
	// content-type header
	encodeJSONError(ctx, err, w)
}

func encodeJSONError(ctx context.Context, err error, w http.ResponseWriter) {
	go apm.CaptureError(ctx, err).Send()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// we can have custom response headers by implementing kithttp.Headerer in
	// our response struct
	if headerer, ok := err.(kithttp.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := errs.ExtractStatusCode(err)
	w.WriteHeader(code)

	var errMessages []string
	switch errT := err.(type) {
	case validator.ValidationErrors:
		errMessages = make([]string, len(errT))
		for i, fieldError := range errT {
			errMessages[i] = fmt.Sprintf(
				"Field validation for '%s' failed on the '%s' tag with value '%v'",
				fieldError.Field(), fieldError.Tag(), fieldError.Value(),
			)
		}
	default:
		errMessages = []string{err.Error()}
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": errMessages,
	})
}
