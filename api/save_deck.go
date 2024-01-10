package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ygo-skc/skc-deck-api/model"
	"github.com/ygo-skc/skc-deck-api/serialization"
)

func submitNewDeckListHandler(res http.ResponseWriter, req *http.Request) {
	var deckList model.DeckList

	if err := json.NewDecoder(req.Body).Decode(&deckList); err != nil {
		log.Println("Error occurred while reading submitNewDeckListHandler request body.")

		res.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(res).Encode(model.APIError{Message: "Body could not be deserialized."})
		return
	}

	log.Printf("Client attempting to submit new deck with name {%s} and with list contents (in base64) {%s}", deckList.Name, deckList.ContentB64)

	// object validation
	if err := deckList.Validate(); err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(res).Encode(err)
		return
	}

	decodedListBytes, _ := base64.StdEncoding.DecodeString(deckList.ContentB64)
	decodedList := string(decodedListBytes) // decoded string of list contents

	var deckListBreakdown model.DeckListBreakdown
	if dlb, err := serialization.DeserializeDeckList(decodedList); err != nil {
		if err.Message == "Could not transform to map" {
			res.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(res).Encode(err)
		} else if err.Message == "Could not access DB" {
			res.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(res).Encode(err)
		}
	} else {
		deckListBreakdown = *dlb
	}

	deckListBreakdown.Sort()

	if err := deckListBreakdown.Validate(); err.Message != "" {
		res.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(res).Encode(err)
		return
	}

	// Adding new deck list, fully validate before insertion
	deckList.ContentB64 = base64.StdEncoding.EncodeToString([]byte(deckListBreakdown.ListStringCleanup()))
	deckList.UniqueCards = deckListBreakdown.CardIDs
	deckList.NumMainDeckCards = deckListBreakdown.NumMainDeckCards
	deckList.NumExtraDeckCards = deckListBreakdown.NumExtraDeckCards
	skcDeckAPIDBInterface.InsertDeckList(deckList)
	json.NewEncoder(res).Encode(model.Success{Message: "Successfully inserted new deck list: " + deckList.Name})
}
