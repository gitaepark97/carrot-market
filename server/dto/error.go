package dto

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/carrot-market/util"
	"github.com/go-playground/validator/v10"
)

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func BindingErrorResponse(err error, obj interface{}, tag string) gin.H {
	if tErr, ok := err.(*json.UnmarshalTypeError); ok {
		err = typeErrorResponse(tErr)
	}

	if vErrs, ok := err.(validator.ValidationErrors); ok {
		err = validatorErrorResponse(vErrs, obj, tag)
	}

	return ErrorResponse(err)
}

func typeErrorResponse(err *json.UnmarshalTypeError) error {
	return util.ErrType(err.Field, err.Type.Name())
}

func validatorErrorResponse(err validator.ValidationErrors, obj interface{}, tag string) error {
	var vErr error

	e := reflect.TypeOf(obj).Elem()
	field, _ := e.FieldByName(err[0].Field())
	tagName, _ := field.Tag.Lookup(tag)

	switch err[0].ActualTag() {
	case "required":
		vErr = util.ErrRequired(tagName)
	case "email":
		vErr = util.ErrEmail(tagName)
	default:
		vErr = err
	}

	return vErr
}
