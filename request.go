package binding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/gocraft/web"
)

var ErrRequestBodyIncomplete = fmt.Errorf("Request body is empty")

func decodeBodyToJSON(ctx interface{}, fieldType reflect.Type, r *web.Request) error {
	t := reflect.ValueOf(ctx)
	if t.Type().Kind() != reflect.Ptr {
		panic("expected pointer to struct")
	}
	t = t.Elem()
	const contextFieldNameForRequest = "RequestJSON"
	saveToField := t.FieldByName(contextFieldNameForRequest)
	if !saveToField.IsValid() {
		panic(fmt.Sprintf("Expected to find field named %q name on the context", contextFieldNameForRequest))
	}

	if !saveToField.CanSet() {
		panic(fmt.Sprintf("Unable to set the value of field named %q on the context", contextFieldNameForRequest))
	}

	newObject := reflect.New(fieldType)
	err := json.NewDecoder(r.Body).Decode(newObject.Elem().Addr().Interface())
	if err != nil {
		if err == io.EOF {
			return ErrRequestBodyIncomplete
		}
		return err
	}

	saveToField.Set(newObject)

	return nil
}

func Request(field interface{}, errorHandlerCustom func(web.ResponseWriter, error)) func(
	interface{}, web.ResponseWriter, *web.Request, web.NextMiddlewareFunc) {

	errorHandler := ErrorHandler
	if errorHandlerCustom != nil {
		errorHandler = errorHandlerCustom
	}

	fieldType := reflect.TypeOf(field)

	return func(ctx interface{}, rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		err := decodeBodyToJSON(ctx, fieldType, r)
		if err != nil {
			errorHandler(rw, err)
			return
		}
		next(rw, r)
	}
}

func ErrorHandler(rw web.ResponseWriter, err error) {
	if rw.Written() {
		panic(fmt.Sprintf("Data already started to be sent to the client and i had an error: %s", err))
	}
	rw.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(rw, "{\"Error\": \"%s\"}", err)
}
