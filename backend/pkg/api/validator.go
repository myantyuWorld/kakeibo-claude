package api

import (
	"backend/pkg/errors"

	playgroundvalidator "github.com/go-playground/validator/v10"
)

type validator struct {
	validate *playgroundvalidator.Validate
}

func newValidator() *validator {
	return &validator{validate: playgroundvalidator.New()}
}

func (v *validator) Validate(i any) error {
	if err := v.validate.Struct(i); err != nil {
		return errors.NewInvalid(err)
	}
	return nil
}
