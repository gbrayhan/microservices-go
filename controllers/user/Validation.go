package user

import (
  "errors"
  "fmt"
  errorModels "github.com/gbrayhan/microservices-go/models/errors"
  "github.com/go-playground/validator/v10"
  "strings"
)

func updateValidation(request map[string]interface{}) (err error) {
  var errorsValidation []string

  for k, v := range request {
    if v == "" {
      errorsValidation = append(errorsValidation, fmt.Sprintf("%s cannot be empty", k))
    }
  }

  validationMap := map[string]string{
    "name":        "omitempty,gt=3,lt=100",
    "description": "omitempty,gt=3,lt=100",
    "ean_code":    "omitempty,gt=3,lt=100",
    "laboratory":  "omitempty,gt=3,lt=100",
  }

  validate := validator.New()
  err = validate.RegisterValidation("update_validation", func(fl validator.FieldLevel) bool {
    m, ok := fl.Field().Interface().(map[string]interface{})
    if !ok {
      return false
    }

    for k, v := range validationMap {
      errValidate := validate.Var(m[k], v)
      if errValidate != nil {
        validatorErr := errValidate.(validator.ValidationErrors)
        errorsValidation = append(errorsValidation, fmt.Sprintf("%s do not satisfy condition %v=%v", k, validatorErr[0].Tag(), validatorErr[0].Param()))
      }
    }

    return true
  })

  if err != nil {
    err = errorModels.NewAppError(err, errorModels.UnknownError)
    return
  }

  err = validate.Var(request, "update_validation")
  if err != nil {
    err = errorModels.NewAppError(err, errorModels.UnknownError)
    return
  }
  if errorsValidation != nil {
    err = errorModels.NewAppError(errors.New(strings.Join(errorsValidation, ", ")), errorModels.ValidationError)
  }
  return
}
