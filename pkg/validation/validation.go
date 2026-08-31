package validation

import (
	"app/product-api/pkg/logs"
	"errors"
	"math"
	"regexp"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
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
	scaled := price * 100
	if math.Abs(scaled-math.Round(scaled)) > 0.00001 {
		return false
	}
	return true
}

func PhoneNumberValidate(f1 validator.FieldLevel) bool {
	phoneNumberRegular, err := regexp.Compile(
		`^(8|\+7)[0-9]{7,10}$`,
	)
	if err != nil {
		logs.AddErrLog(log.Fields{
			"Error": err,
		}, "Ошибка валидации регулярного выражения")
		return false
	}

	isValid := phoneNumberRegular.MatchString(f1.Field().String())
	if !isValid {
		return false
	}

	return true
}
