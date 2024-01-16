package validation

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/ygo-skc/skc-deck-api/model"
)

// Add custom validators to handle validation scenarios not supported out of the box.
func configureCustomValidators() {
	V.RegisterValidation(deckListNameValidator, func(fl validator.FieldLevel) bool {
		return len(deckListNameRegex.FindAllString(fl.Field().String(), -1)) > 0
	})

	V.RegisterValidation(deckMascotsValidator, func(fl validator.FieldLevel) bool {
		mascots := fl.Field().Interface().(model.CardIDs)

		for ind, mascot := range mascots {
			if ind == 3 { // size constraint fails
				log.Println("Deck Mascot array failed size constraint.")
				return false
			} else if len(cardIDRegex.FindAllString(mascot, -1)) == 0 { // regex constraint
				log.Println("Deck Mascot ID not in proper format.")
				return false
			}
		}

		return true
	})
}

// validate deck list
func Validate(dl model.DeckList) *ValidationErrors {
	if err := V.Struct(dl); err != nil {
		return handleValidationErrors(err.(validator.ValidationErrors))
	} else {
		return nil
	}
}
