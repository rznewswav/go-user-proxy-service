package config

import "errors"

var ErrReflectFieldNotSetable = errors.New(
	"reflect field is not settable",
)

var ErrCannotCast = errors.New(
	"environment value cannot be casted to the expected type",
)

var ErrNotImplemented = errors.New(
	"reflect field type is not implemented yet",
)

const envTagName = "env"
const isRequiredTagName = "required"
const defaultValueTagName = "default"
const printDebugValueTagName = "printDebug"

var falsyString = map[string]bool{
	"false": false,
	"off":   false,
	"0":     false,
}
