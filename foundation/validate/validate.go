// Package validate contains the support for validating models.
package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

func init() {

	// Instantiate a validator.
	validate = validator.New()
	validate.RegisterValidation("currency", validCurrency)
	validate.RegisterValidation("password", validPassword)
	validate.RegisterValidation("username", validUsername)
	validate.RegisterValidation("accountname", validAccountName)

	// Create a translator for english so the error messages are
	// more human-readable than technical.
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// Register the english error messages for use.
	en_translations.RegisterDefaultTranslations(validate, translator)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// require currency to be one of: USD, EUR, CAD
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		switch currency {
		case USD, EUR, CAD:
			return true
		}
		return false
	}
	return false
}

// require password to be 6 ~ 20 characters long
var validPassword validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if password, ok := fieldLevel.Field().Interface().(string); ok {
		if len(password) < 6 || len(password) > 20 {
			return false
		}
		return true
	}
	return false
}

// require username to be 3 ~ 20 characters long
var validUsername validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if username, ok := fieldLevel.Field().Interface().(string); ok {
		if len(username) < 3 || len(username) > 20 {
			return false
		}
		return true
	}
	return false

}

// require account name to be less than 30 characters
var validAccountName validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if name, ok := fieldLevel.Field().Interface().(string); ok {
		return len(name) <= 30
	}
	return false
}

// Check validates the provided model against it's declared tags.
func Check(val any) error {
	if err := validate.Struct(val); err != nil {
		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Err:   verror.Translate(translator),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}
