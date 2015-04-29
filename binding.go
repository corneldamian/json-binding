package binding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gocraft/web"
)

func decodeBodyToJSON(ctx interface{}, fieldType reflect.Type, r *web.Request) error {
	t := reflect.ValueOf(ctx)
	if t.Type().Kind() != reflect.Ptr {
		panic("expected pointer to struct")
	}
	t = t.Elem()
	saveToField := t.FieldByName("BodyJSON")
	if !saveToField.IsValid() {
		panic("Expected to find BodyJSON field name on the context")
	}

	if !saveToField.CanSet() {
		panic("Unable to set BodyJSON field value on the context")
	}

	newObject := reflect.New(fieldType)
	err := json.NewDecoder(r.Body).Decode(newObject.Elem().Addr().Interface())
	if err != nil {
		return err
	}

	saveToField.Set(newObject)

	return nil
}

func Bind(field interface{}, errorHandlerCustom func(web.ResponseWriter, error)) func(
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
	rw.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(rw, "{\"error\": \"%s\"", err)
}
