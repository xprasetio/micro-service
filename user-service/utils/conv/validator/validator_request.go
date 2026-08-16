package validator

import (
	"errors"
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatalf("[New Validator-1] translator not found")
	}
	validate := validator.New()
	return &Validator{
		Validator:  validate,
		Translator: trans,
	}
}

func (v *Validator) Validate(i interface{}) error {
	err := v.Validator.Struct(i)
	if err != nil {
		object, _ := err.(validator.ValidationErrors)
		for _, e := range object {
			log.Printf("[Validate-1] %s:%s", e.Field(), e.Translate(v.Translator))
			return errors.New(e.Translate(v.Translator))
		}
	}
	return nil
}
