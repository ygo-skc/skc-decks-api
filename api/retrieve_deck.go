package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ygo-skc/skc-deck-api/model"
	"github.com/ygo-skc/skc-deck-api/serialization"
)

func getDeckListHandler(res http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	deckID := pathVars["deckID"]
	log.Println("Getting content for deck w/ ID:", deckID)

	var deckList *model.DeckList
	var err *model.APIError
	if deckList, err = skcDeckAPIDBInterface.GetDeckList(deckID); err != nil {
		err.HandleServerResponse(res)
		return
	}

	decodedListBytes, _ := base64.StdEncoding.DecodeString(deckList.ContentB64)
	decodedList := string(decodedListBytes) // decoded string of list contents

	var deckListBreakdown *model.DeckListBreakdown
	if deckListBreakdown, err = serialization.DeserializeDeckList(decodedList); err != nil {
		err.HandleServerResponse(res)
		return
	}

	mainDeckContent := make([]model.Content, 0, len(deckListBreakdown.MainDeck))
	for _, card := range deckListBreakdown.MainDeck {
		mainDeckContent = append(mainDeckContent, model.Content{Card: card, Quantity: deckListBreakdown.CardQuantity[card.CardID]})
	}
	deckList.MainDeck = &mainDeckContent

	extraDeck := make([]model.Content, 0, len(deckListBreakdown.ExtraDeck))
	for _, card := range deckListBreakdown.ExtraDeck {
		extraDeck = append(extraDeck, model.Content{Card: card, Quantity: deckListBreakdown.CardQuantity[card.CardID]})
	}
	deckList.ExtraDeck = &extraDeck

	log.Printf("Successfully retrieved deck list. Name {%s} and encoded deck list content {%s}. This deck list has {%d} main deck cards and {%d} extra deck cards.", deckList.Name, deckList.ContentB64, deckList.NumMainDeckCards, deckList.NumExtraDeckCards)
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(deckList)
}

func getDecksFeaturingCardHandler(res http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	cardID := pathVars["cardID"]
	log.Printf("Getting decks that use card w/ ID: %s", cardID)

	suggestedDecks := model.SuggestedDecks{}

	suggestedDecks.FeaturedIn, _ = skcDeckAPIDBInterface.GetDecksThatFeatureCards([]string{cardID})

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(suggestedDecks)
}
