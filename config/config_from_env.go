package config

import (
	"errors"
	"fmt"
	"github.com/IntelXLabs-LLC/go-common/utils"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var defaultSettings DefaultSettings

// GetCurrentEnvironment returns the current environment based on the default settings.
func GetCurrentEnvironment() Environment {
	return defaultSettings.Environment
}

// LoadConfig loads the configuration for the given configStruct by setting appropriate values based on the current environment.
// It imports environment variables from a `.env` file in the current directory and checks if the environment variable "ENVIRONMENT" matches any of the environment keys (UAT, Staging
func LoadConfig(configStruct interface{}) error {
	importEnv()
	val := viper.GetString("ENVIRONMENT")
	val = strings.ToLower(val)
	switch val {
	case "uat":
		defaultSettings.Environment = UAT
	case "stg":
		defaultSettings.Environment = Staging
	case "prod":
		defaultSettings.Environment = Production
	default:
		defaultSettings.Environment = Development
	}
	return loadConfigRecursive(reflect.ValueOf(configStruct).Elem(), "", false)
}

// loadConfigRecursive traverses a struct recursively and loads configuration values from environment variables or default values.
// It takes in the structValue to be populated, prefix for potential environment variable prefixes, and optionalFields to specify which fields are optional.
// It returns an error if any required fields are missing.
func loadConfigRecursive(structValue reflect.Value, prefix string, isOptional bool, optionalFields ...string) error {
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		if field.Type.Kind() == reflect.Struct {
			prefix := field.Tag.Get("prefix")
			isOptional := field.Tag.Get("optional") == "true"
			optionalFieldsStr := field.Tag.Get("optional_fields")
			optionalFields := make([]string, 0)
			if optionalFieldsStr != "" {
				optionalFields = strings.Split(optionalFieldsStr, ";")
			}
			if err := loadConfigRecursive(fieldValue, prefix, isOptional, optionalFields...); err != nil {
				return err
			}
		}

		envName := field.Tag.Get("config")
		defaultValue := field.Tag.Get("default")
		isOptional := field.Tag.Get("optional") == "true" || utils.IsInArray(envName, optionalFields)

		if envName == "" {
			continue
		} else if prefix != "" {
			envName = fmt.Sprintf("%s_%s", prefix, envName)
		}

		if value := viper.GetString(envName); value != "" {
			setFieldValue(envName, fieldValue, value, field.Type.Kind(), false)
		} else if defaultValue != "" {
			setFieldValue(envName, fieldValue, defaultValue, field.Type.Kind(), true)
		} else {
			if isOptional {
				continue
			}
			return fmt.Errorf("missing value for %s", field.Name)
		}
	}

	return nil
}

// setFieldValue sets the value of a field in a struct based on the provided parameters.
//
// Parameters:
// - envName: the environment variable name
// - fieldValue: the reflect.Value of the field to be set
// - value: the value to set for the field
// - kind: the reflect.Kind of the field
// - parseDefault: flag indicating whether to parse the default value or not
//
// Note: This function handles different kinds of fields including String, Integers, Booleans.
// If parseDefault is true, the function will set the field value to the provided default value.
// If parseDefault is false, the function will get the value from the environment variable.
//
// This function panics if the field type is not implemented.
//
// Example usage:
//
//	fieldType := reflect.TypeOf(MyStruct{}.MyField)
//	fieldValue := reflect.New(fieldType).Elem()
//	setFieldValue("ENV_VAR_NAME", fieldValue, "default value", fieldType.Kind(), true)
//
// For additional information, refer to the following helper functions:
// - parseInteger: to parse integer values
// - parseBool: to parse boolean values
// - utils.IsInArray: to check if a value is present in an array.
func setFieldValue(envName string, fieldValue reflect.Value, value string, kind reflect.Kind, parseDefault bool) {
	switch kind {
	case reflect.String:
		if !parseDefault {
			fieldValue.SetString(viper.GetString(envName))
		} else {
			fieldValue.SetString(value)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !parseDefault {
			fieldValue.SetInt(viper.GetInt64(envName))
		} else {
			if intValue, err := parseInteger(value, fieldValue.Type().Bits()); err == nil {
				fieldValue.SetInt(intValue)
			}
		}
	case reflect.Bool:
		if !parseDefault {
			fieldValue.SetBool(viper.GetBool(envName))
		} else {
			if boolValue, err := parseBool(value); err == nil {
				fieldValue.SetBool(boolValue)
			}
		}
	default:
		panic("unhandled case: type is not implemented")
	}
}

// parseInteger converts a string value to a signed integer of the specified number of bits.
func parseInteger(value string, bits int) (int64, error) {
	return strconv.ParseInt(value, 10, bits)
}

// parseBool parses the given string value and returns a boolean value.
func parseBool(value string) (bool, error) {
	return strconv.ParseBool(value)
}

// importEnv sets up the configuration for reading environment variables.
// It configures viper to look for a config file named ".env" in the current directory.
// It then adds the current directory as a config path and sets the config type to "env".
// The function also checks if environment variables match any of the existing keys and loads them using viper.AutomaticEnv().
// If there is an error reading the config file, it checks if the error is of type viper.ConfigFileNotFoundError and ignores it.
// Usage example:
//
//	importEnv()
func importEnv() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// checks if environment variables match any of the existing keys and loads them.
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Config file not found ignoring error
		}
	}

}
