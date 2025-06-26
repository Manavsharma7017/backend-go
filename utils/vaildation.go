package utils

import (
	"backend/models"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		role := fl.Field().String()

		// Get the struct type
		structType := fl.Parent().Type()

		switch structType {
		case reflect.TypeOf(models.User{}):
			// User must always have RoleUSER only
			return role == string(models.RoleUSER)
		case reflect.TypeOf(models.AdminUser{}):
			// AdminUser can be RoleAdmin or RoleEditor
			return role == string(models.RoleAdmin) || role == string(models.RoleEditor)
		default:
			return false
		}
	})
}

// ValidateStruct validates any input struct and prints the field errors
func ValidateStruct(input interface{}) error {
	err := validate.Struct(input)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			fmt.Printf("Validation failed for field '%s': %s\n", fieldErr.Field(), fieldErr.ActualTag())
		}
		return err
	}
	return nil
}
