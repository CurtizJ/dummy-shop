package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func NewConfigFromEnv(config interface{}) error {
	value := reflect.ValueOf(config).Elem()
	valueType := reflect.TypeOf(config).Elem()

	for i := 0; i < valueType.NumField(); i++ {
		typeField := valueType.Field(i)
		envName, exists := typeField.Tag.Lookup("env")
		if !exists {
			continue
		}

		defaultValue, defaultExists := typeField.Tag.Lookup("default")
		envValue, envExists := os.LookupEnv(envName)

		if !envExists {
			if !defaultExists {
				return fmt.Errorf("Neither env '%s' exists neither default is specified", envName)
			}

			envValue = defaultValue
		}

		valueField := value.Field(i)
		kind := valueField.Kind()

		if kind >= reflect.Int && kind <= reflect.Int64 {
			intValue, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return err
			}
			valueField.SetInt(intValue)
		} else if kind >= reflect.Uint && kind <= reflect.Uint64 {
			uintValue, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return err
			}
			valueField.SetUint(uintValue)
		} else if kind == reflect.String {
			valueField.SetString(envValue)
		} else {
			return fmt.Errorf("Unsupported type %T for field %s", valueField.Type, typeField.Name)
		}
	}

	return nil
}
