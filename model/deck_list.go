package model

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ygo-skc/skc-deck-api/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SuggestedDecks struct {
	FeaturedIn *[]DeckList `json:"featuredIn"`
}

type DeckList struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name              string             `bson:"name" json:"name" validate:"required,decklistname"`
	ContentB64        string             `bson:"content" json:"listContent" validate:"required,base64"`
	VideoUrl          string             `bson:"videoUrl" json:"videoUrl" validate:"omitempty,url"`
	UniqueCards       []string           `bson:"uniqueCards" json:"uniqueCards" validate:"omitempty"`
	DeckMascots       []string           `bson:"deckMascots" json:"deckMascots" validate:"omitempty,deckmascots"`
	NumMainDeckCards  int                `bson:"numMainDeckCards" json:"numMainDeckCards"`
	NumExtraDeckCards int                `bson:"numExtraDeckCards" json:"numExtraDeckCards"`
	Tags              []string           `bson:"tags" json:"tags" validate:"required"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
	MainDeck          *[]Content         `bson:"mainDeck,omitempty" json:"mainDeck,omitempty"`
	ExtraDeck         *[]Content         `bson:"extraDeck,omitempty" json:"extraDeck,omitempty"`
}

type Content struct {
	Quantity int  `bson:"omitempty" json:"quantity"`
	Card     Card `bson:"omitempty" json:"card"`
}

type DeckListBreakdown struct {
	CardQuantity      map[string]int
	CardIDs           []string
	InvalidIDs        []string
	AllCards          CardDataMap
	MainDeck          []Card
	ExtraDeck         []Card
	NumMainDeckCards  int
	NumExtraDeckCards int
}

// validate and handle validation error messages
func (dl DeckList) Validate() *validation.ValidationErrors {
	if err := validation.V.Struct(dl); err != nil {
		return validation.HandleValidationErrors(err.(validator.ValidationErrors))
	} else {
		return nil
	}
}

func (dlb *DeckListBreakdown) Sort() {
	dlb.MainDeck = []Card{}
	dlb.ExtraDeck = []Card{}
	numMainDeckCards := 0
	numExtraDeckCards := 0

	for _, cardID := range dlb.CardIDs {
		if _, isPresent := dlb.AllCards[cardID]; isPresent {
			if dlb.AllCards[cardID].IsExtraDeckMonster() {
				dlb.ExtraDeck = append(dlb.ExtraDeck, dlb.AllCards[cardID])
				numExtraDeckCards += dlb.CardQuantity[cardID]
			} else {
				dlb.MainDeck = append(dlb.MainDeck, dlb.AllCards[cardID])
				numMainDeckCards += dlb.CardQuantity[cardID]
			}
		}
	}

	dlb.NumMainDeckCards = numMainDeckCards
	dlb.NumExtraDeckCards = numExtraDeckCards

	sortDeckUsingName(&dlb.MainDeck)
	sortDeckUsingName(&dlb.ExtraDeck)
}

func sortDeckUsingName(cards *[]Card) {
	sort.SliceStable(*cards, func(i, j int) bool {
		return (*cards)[i].CardName < (*cards)[j].CardName
	})
}

func (dlb DeckListBreakdown) ListStringCleanup() string {
	formattedDLS := "Main Deck\n"

	for _, card := range dlb.MainDeck {
		formattedDLS += formattedLine(card, dlb.CardQuantity[card.CardID])
	}

	formattedDLS += "\nExtra Deck\n"

	for _, card := range dlb.ExtraDeck {
		formattedDLS += formattedLine(card, dlb.CardQuantity[card.CardID])
	}

	return formattedDLS
}

func formattedLine(card Card, quantity int) string {
	return fmt.Sprintf("%dx%s|%s\n", quantity, card.CardID, card.CardName)
}

func (dlb DeckListBreakdown) Validate() *APIError {
	var msg = ""

	if len(dlb.InvalidIDs) > 0 {
		msg = fmt.Sprintf("Deck list contains card(s) that were not found in skc DB. All cards not found in DB: %v", dlb.InvalidIDs)
	}

	// validate extra deck has correct number of cards
	if dlb.NumExtraDeckCards > 15 {
		msg = fmt.Sprintf("Extra deck cannot contain more than 15 cards. Current deck contains %d extra deck cards.", dlb.NumExtraDeckCards)
	}

	// validate main deck has correct number of cards
	if dlb.NumMainDeckCards < 40 || dlb.NumMainDeckCards > 60 {
		msg = fmt.Sprintf("Main deck cannot contain less than 40 cards and no more than 60 cards. Current deck contains %d main deck cards.", dlb.NumMainDeckCards)
	}

	if msg != "" {
		log.Println(msg)
		return &APIError{Message: msg, StatusCode: http.StatusBadRequest}
	} else {
		return nil
	}
}
