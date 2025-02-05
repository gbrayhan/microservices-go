package user

import (
	"errors"
	"fmt"
	"strings"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/go-playground/validator/v10"
)

func updateValidation(request map[string]any) error {
	var errorsValidation []string
	for k, v := range request {
		if v == "" {
			errorsValidation = append(errorsValidation, fmt.Sprintf("%s cannot be empty", k))
		}
	}

	validationMap := map[string]string{
		"user_name": "omitempty,gt=3,lt=100",
		"email":     "omitempty,email",
		"firstName": "omitempty,gt=1,lt=100",
		"lastName":  "omitempty,gt=1,lt=100",
	}

	validate := validator.New()
	err := validate.RegisterValidation("update_validation", func(fl validator.FieldLevel) bool {
		m, ok := fl.Field().Interface().(map[string]any)
		if !ok {
			return false
		}
		for k, rule := range validationMap {
			if val, exists := m[k]; exists {
				errValidate := validate.Var(val, rule)
				if errValidate != nil {
					validatorErr := errValidate.(validator.ValidationErrors)
					errorsValidation = append(
						errorsValidation,
						fmt.Sprintf("%s does not satisfy condition %v=%v", k, validatorErr[0].Tag(), validatorErr[0].Param()),
					)
				}
			}
		}
		return true
	})
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.UnknownError)
	}

	err = validate.Var(request, "update_validation")
	if err != nil {
		return domainErrors.NewAppError(err, domainErrors.UnknownError)
	}
	if len(errorsValidation) > 0 {
		return domainErrors.NewAppError(errors.New(strings.Join(errorsValidation, ", ")), domainErrors.ValidationError)
	}
	return nil
}
