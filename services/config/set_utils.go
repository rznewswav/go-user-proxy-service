package config

import (
	"fmt"
	"reflect"
	"strconv"
)

func SetInt(
	reflectValue reflect.Value,
	strValue string,
) error {
	i, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return fmt.Errorf(
			"%s: %s to type int64",
			ErrCannotCast.Error(),
			strValue,
		)
	}

	reflectValue.SetInt(i)
	return nil
}

func SetBool(
	reflectValue reflect.Value,
	strValue string,
) error {
	_, isFalse := falsyString[strValue]
	reflectValue.SetBool(!isFalse)

	return nil
}

func SetString(
	reflectValue reflect.Value,
	strValue string,
) error {
	reflectValue.SetString(strValue)

	return nil
}

func SetReflectValue(
	reflectValue reflect.Value,
	strValue string,
) error {
	if !reflectValue.CanSet() {
		return ErrReflectFieldNotSetable
	}
	reflectValueKind := reflectValue.Kind()
	switch reflectValueKind {
	case reflect.Int64:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int:
		return SetInt(reflectValue, strValue)
	case reflect.Bool:
		return SetBool(reflectValue, strValue)
	case reflect.String:
		return SetString(reflectValue, strValue)
	}
	return fmt.Errorf(
		"%w: type %s",
		ErrNotImplemented,
		reflectValueKind.String(),
	)
}
