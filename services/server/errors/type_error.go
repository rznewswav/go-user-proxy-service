package errors

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
	switch t.Type().Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		translated, translationError := ut.T("number", t.Field())
		if translationError != nil {
			goto END
		}
		return translated
	case reflect.Bool:
		translated, translationError := ut.T("boolean", t.Field())
		if translationError != nil {
			goto END
		}
		return translated
	case reflect.String:
		translated, translationError := ut.T("string", t.Field())
		if translationError != nil {
			goto END
		}
		return translated
	}

END:
	return fmt.Sprintf("%s must be of type %s", t.Field(), t.Type().Kind())
}

func (t TypeError) Error() string {
	return fmt.Sprintf(
		"expected data type \"%s\" for \"%s\" but got type \"%s\" instead",
		t.Type().Kind(),
		t.Field(),
		t.CurrentValueType,
	)
}
