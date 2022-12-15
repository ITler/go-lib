package misc

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	uni *ut.UniversalTranslator
	// Translator provides a [ut.Translator] instance
	Translator ut.Translator
	// StructValidator provides a [validator.Validate] instance
	StructValidator *validator.Validate
)

// CanBeValidated describes a type that can be validated
type CanBeValidated interface {
	Validate() error
}

// Validator encapsulates information about how to validate a certain struct field
type Validator struct {
	Field     string
	Message   string
	Validator func(fl validator.FieldLevel) bool
}

// ValidatorsCreatorFunc describes a signature for functions that should be able
// to create [Validator]s
type ValidatorsCreatorFunc func(in CanBeValidated) []Validator

// RegisterValidatorsDefault allows registering validators and are held within this library
func RegisterValidatorsDefault(validators []Validator) {
	RegisterValidators(validators, StructValidator, Translator)
}

// RegisterValidators helps to register [Validator]s
// with provided translations and validations
func RegisterValidators(validators []Validator, v *validator.Validate, trans ut.Translator) {
	for _, val := range validators {
		_ = v.RegisterValidation(val.Field, val.Validator)
		_ = v.RegisterTranslation(val.Field, trans,
			func(ut ut.Translator) error {
				return ut.Add(val.Field, val.Message, true) // see universal-translator for details
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(val.Field, fe.Field())
				return t
			})
	}
}

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	translator := en.New()
	uni = ut.New(translator, translator)

	Translator, _ = uni.GetTranslator("en")

	StructValidator = validator.New()

	if err := en_translations.RegisterDefaultTranslations(StructValidator, Translator); err != nil {
		log.Fatal().Stack().Err(err).Msg("")
	}
}
