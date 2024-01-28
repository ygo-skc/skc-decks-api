package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/ygo-skc/skc-deck-api/model"
)

// validate deck list
func Validate(dl model.DeckList) *ValidationErrors {
	if err := V.Struct(dl); err != nil {
		return handleValidationErrors(err.(validator.ValidationErrors))
	} else {
		return nil
	}
}
