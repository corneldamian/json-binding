package binding

import (
	"encoding/json"
	"reflect"

	"github.com/gocraft/web"
)

type response struct {
	Error     *string     `json:",omitempty"`
	ErrorCode *int        `json:",omitempty"`
	Success   *string     `json:",omitempty"`
	Data      interface{} `json:",omitempty"`
}

var ok = "ok"

func SuccessResponse(data interface{}) *response {
	if s, ok := data.(string); ok {
		return &response{
			Success: &s,
		}
	}

	return &response{
		Success: &ok,
		Data:    data,
	}
}

func ErrorResponse(err string, code int) *response {
	return &response{
		Error:     &err,
		ErrorCode: &code,
	}
}

func Response(errorHandlerCustom func(web.ResponseWriter, error)) func(
	interface{}, web.ResponseWriter, *web.Request, web.NextMiddlewareFunc) {

	errorHandler := ErrorHandler
	if errorHandlerCustom != nil {
		errorHandler = errorHandlerCustom
	}

	return func(ctx interface{}, rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		next(rw, r)

		t := reflect.ValueOf(ctx)
		if t.Type().Kind() != reflect.Ptr {
			panic("expected pointer to struct")
		}
		t = t.Elem()

		var (
			data          []byte
			err           error
			responseField = t.FieldByName("ResponseJSON")
		)

		if responseField.IsValid() && !responseField.IsNil() {
			if responseField.Kind() == reflect.String {
				data = []byte(responseField.String())
			} else {
				data, err = json.Marshal(responseField.Addr().Interface())
				if err != nil {
					errorHandler(rw, err)
					return
				}
			}

		}

		if rw.StatusCode() == 0 {
			statusCodeField := t.FieldByName("ResponseStatus")
			if statusCodeField.IsValid() {
				statusCode := statusCodeField.Int()
				if statusCode > 0 {
					rw.WriteHeader(int(statusCode))
				}
			}
		}

		if len(data) > 0 {
			_, err = rw.Write(data)
			if err != nil {
				rw.Write([]byte("{\"Error\": \"writing error\"")) //try to send error ???
				panic("I was unable to write data to the client: " + err.Error())
			}
		}
	}
}
