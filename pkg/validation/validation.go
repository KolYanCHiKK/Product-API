package validation

import (
	"errors"
	"math"

	"github.com/go-playground/validator/v10"
)

func ValidateBody(body any, validateMap map[string]string, registerFunc func() (*validator.Validate, error)) []string {
	var validateErr []string
	validate, err := registerFunc()
	if err != nil {
		validateErr = append(validateErr, err.Error())
		return validateErr
	}
	err = validate.Struct(body)

	var errSlice validator.ValidationErrors
	if !errors.As(err, &errSlice) {
		return nil
	}

	for _, value := range errSlice {
		key := value.Field() + "." + value.Tag()
		errMess := validateMap[key]
		validateErr = append(validateErr, errMess)
	}
	return validateErr
}

func ScaleValidate(f1 validator.FieldLevel) bool {
	price := f1.Field().Float()
	if math.Round(price*100) != price*100 {
		return false
	}
	return true
}
