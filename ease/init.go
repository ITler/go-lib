package ease

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
