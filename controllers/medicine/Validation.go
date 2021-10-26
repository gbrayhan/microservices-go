package medicine

import (
  errorModels "github.com/gbrayhan/microservices-go/models/errors"
  "github.com/go-playground/validator/v10"
)

func updateValidation(request map[string]interface{}) (err error) {
  validationMap := map[string]string{
    "name":        "omitempty,min=1,max=100",
    "description": "omitempty,min=1,max=100",
    "ean_code":    "omitempty,min=1,max=100",
    "laboratory":  "omitempty,min=1,max=100",
  }

  validate := validator.New()
  err = validate.RegisterValidation("update_validation", func(fl validator.FieldLevel) bool {
    m, ok := fl.Field().Interface().(map[string]interface{})
    if !ok {
      return false
    }

    for k, v := range validationMap {
      if validate.Var(m[k], v) != nil {
        return false
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

  return
}
