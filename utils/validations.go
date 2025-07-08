package utils

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(request interface{}) []ErrorResponse {
	var errors []ErrorResponse
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, element)
		}
	}
	return errors
}

func InitValidators() {
	validate = validator.New()
	_ = validate.RegisterValidation("full_name", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return len(s) <= 64
	})
	_ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return validatePassword(fl.Field().String())
	})

	err := validate.RegisterValidation("nest_command", func(fl validator.FieldLevel) bool {
		return validateGoogleNestCommand(fl)
	})
	if err != nil {
		panic(err)
	}

}

func validateGoogleNestCommand(fl validator.FieldLevel) bool {
	var availableModes = []string{
		"sdm.devices.commands.ThermostatEco.SetMode",
		"sdm.devices.commands.ThermostatMode.SetMode",
		"sdm.devices.commands.ThermostatTemperatureSetpoint.SetHeat",
		"sdm.devices.commands.ThermostatTemperatureSetpoint.SetCool",
		"sdm.devices.commands.ThermostatTemperatureSetpoint.SetRange",
                "sdm.devices.commands.Fan.SetTimer",
	}

	for _, mode := range availableModes {
		if fl.Field().String() == mode {
			return true
		}
	}
	return false
}

func validatePassword(password string) bool {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)
	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber
}

func init() {
	InitValidators()
}
