package model

type QuotedToken = string

// identifier arrays

type CardIDs []string
type CardNames []string
type IdentifierSlice interface {
	CardIDs | CardNames
}

// data maps

type CardDataMap map[string]Card

// finds all card IDs not found in CardDataMap keys
func (cardData CardDataMap) FindMissingIDs(cardIDs CardIDs) CardIDs {
	missingIDs := make(CardIDs, 0)

	for _, cardID := range cardIDs {
		if _, containsKey := cardData[cardID]; !containsKey {
			missingIDs = append(missingIDs, cardID)
		}
	}

	return missingIDs
}

// finds all card IDs not found in CardDataMap keys
func (cardData CardDataMap) FindMissingNames(cardNames CardNames) CardNames {
	missingNames := make(CardNames, 0)

	for _, cardName := range cardNames {
		if _, containsKey := cardData[cardName]; !containsKey {
			missingNames = append(missingNames, cardName)
		}
	}

	return missingNames
}

type ResourceDataMap interface {
	CardDataMap
}

// data types that contain many resources of the same data type

type BatchCardIDs struct {
	CardIDs CardIDs `json:"cardIDs" validate:"required,ygocardids"`
}

type BatchCardData[IS IdentifierSlice] struct {
	CardInfo         CardDataMap `json:"cardInfo"`
	UnknownResources IS          `json:"unknownResources"`
}

type BatchData[IS IdentifierSlice] interface {
	BatchCardData[IS]
}
