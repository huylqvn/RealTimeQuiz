package validator

import (
	"quizserver/src/singleton"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validatorSingleton = singleton.NewSingleton(func() *Validator { return NewValidator() }, true)
)

type Validator struct {
	validator *validator.Validate
	trans     ut.Translator
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	v := &Validator{validator: validate, trans: trans}
	v.init()
	return v
}

func GetValidator() *Validator {
	return validatorSingleton.Get()
}

func (v *Validator) init() {
	v.validator.RegisterTagNameFunc(getJsonTagName)
}

func (v Validator) Validate(data interface{}) error {
	err := v.validator.Struct(data)
	if err != nil {
		return err
		// validationErrors := []error{}
		// for _, err := range errs.(validator.ValidationErrors) {
		// 	// In this case data object is actually holding the User struct
		// 	var elem ValidateError

		// 	elem.FailedField = err.Field() // Export struct field name
		// 	elem.Tag = err.Tag()           // Export struct tag
		// 	elem.Value = err.Value()       // Export field value

		// 	validationErrors = append(validationErrors, elem)
		// }
	}

	return nil
}

func getJsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}
