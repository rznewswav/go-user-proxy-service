package controllers

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"reflect"
)

type TypeError struct {
	FieldName        string
	CurrentValueType string
	TypeOf           reflect.Type
}

func (t TypeError) Tag() string {
	return "type"
}

func (t TypeError) ActualTag() string {
	return ""
}

func (t TypeError) Namespace() string {
	return ""
}

func (t TypeError) StructNamespace() string {
	return ""
}

func (t TypeError) Field() string {
	return t.FieldName
}

func (t TypeError) StructField() string {
	return t.FieldName

}

func (t TypeError) Value() interface{} {
	return ""
}

func (t TypeError) Param() string {
	return ""
}

func (t TypeError) Kind() reflect.Kind {
	return reflect.Interface
}

func (t TypeError) Type() reflect.Type {
	return t.TypeOf
}

func (t TypeError) Translate(ut ut.Translator) string {
	return ""
}

func (t TypeError) Error() string {
	return fmt.Sprintf(
		"expected data type \"%s\" for \"%s\" but got type \"%s\" instead",
		t.Type().Kind(),
		t.Field(),
		t.CurrentValueType,
	)
}
