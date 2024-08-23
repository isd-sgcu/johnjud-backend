package validator

import (
	"errors"
	"net/http"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/isd-sgcu/johnjud-backend/internal/dto"
	"github.com/rs/zerolog/log"
)

type IDtoValidator interface {
	Validate(interface{}) []*dto.BadReqErrResponse
}

type DtoValidator struct {
	v     *validator.Validate
	trans ut.Translator
}

func (v *DtoValidator) Validate(in interface{}) []*dto.BadReqErrResponse {
	err := v.v.Struct(in)

	var errors []*dto.BadReqErrResponse
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			element := dto.BadReqErrResponse{
				Message:     e.Translate(v.trans),
				FailedField: e.StructField(),
				Value:       e.Value(),
			}

			log.Error().
				Str("module", "validate").
				Int("status", http.StatusBadRequest).
				Interface("error", element).
				Msg("Validate failed")

			errors = append(errors, &element)
		}
	}
	return errors
}

func NewValidator() (*DtoValidator, error) {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, errors.New("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, err
	}

	return &DtoValidator{
		v:     v,
		trans: trans,
	}, nil
}

func NewIValidator() (IDtoValidator, error) {
	return NewValidator()
}
