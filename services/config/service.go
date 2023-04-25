package config

import (
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/joho/godotenv"
)

type EnvItem struct {
	name         string
	value        string
	err          error
	reflectValue reflect.Value
	usingDefault bool
	printOnDebug bool
}

func init() {
	cwd := os.Getenv("CWD")
	if len(cwd) > 0 {
		godotenv.Load(path.Join(cwd, ".env"))
	} else {
		godotenv.Load()
	}
}

func Build[T any](constructable T) (T, error) {
	typeof := reflect.TypeOf(constructable)
	configPointerToStruct := reflect.ValueOf(&constructable).
		Elem()

	envItems := make([]EnvItem, 0)

	for j := 0; j < typeof.NumField(); j++ {
		envItems = append([]EnvItem{{}}, envItems...)
		envItem := &envItems[0]

		configValueField := typeof.Field(j)
		configValueFieldStruct := configPointerToStruct.FieldByName(
			configValueField.Name,
		)

		// Get the field tag value
		env := configValueField.Tag.Get(envTagName)
		isRequired := configValueField.Tag.Get(
			isRequiredTagName,
		)
		defaultValue := configValueField.Tag.Get(
			defaultValueTagName,
		)

		currentValue := os.Getenv(env)

		if len(env) == 0 {
			envItem.err = fmt.Errorf(
				"env is not set in struct tag",
			)
			continue
		}

		var valueToSet = currentValue

		if len(currentValue) == 0 {
			valueToSet = defaultValue
			envItem.usingDefault = true
		}

		envItem.name = configValueField.Name
		envItem.value = valueToSet
		envItem.printOnDebug = configValueField.Tag.Get(
			printDebugValueTagName,
		) == "true"

		if len(valueToSet) == 0 && isRequired == "true" {
			envItem.err = fmt.Errorf(
				"%s is not set in .env",
				env,
			)
			continue
		}

		setReflectError := SetReflectValue(
			configValueFieldStruct,
			valueToSet,
		)

		if setReflectError != nil {
			envItem.err = setReflectError
		}
	}
	var errorOrNil error

	wrapToErrorClass := func(err error) {
		if errorOrNil == nil {
			errorOrNil = fmt.Errorf(
				"error on setting app configuration:\n    %s",
				err,
			)
		} else {
			errorOrNil = fmt.Errorf("%s\n    %s", errorOrNil.Error(), err)
		}
	}

	for _, env := range envItems {
		if env.err != nil {
			errorToWrap := fmt.Errorf(
				"error on setting config field %s: %w",
				env.name,
				env.err,
			)
			wrapToErrorClass(errorToWrap)
			continue
		}

		SetReflectValue(env.reflectValue, env.value)

		if env.printOnDebug {
			fmt.Printf(
				"Using env %s: %s\n",
				env.name,
				env.value,
			)
		}
	}

	return constructable, errorOrNil
}

func QuietBuild[T any](constructable T) T {
	constructable, _ = Build(constructable)
	return constructable
}
