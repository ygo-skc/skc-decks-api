package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ygo-skc/skc-deck-api/io"
	"github.com/ygo-skc/skc-deck-api/model"
	cModel "github.com/ygo-skc/skc-go/common/model"
	cUtil "github.com/ygo-skc/skc-go/common/util"
)

func getDeckListHandler(res http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	deckID := pathVars["deckID"]

	logger, ctx := cUtil.NewRequestSetup(context.Background(), "retrieve deck list", slog.String("deckID", deckID))
	logger.Info(fmt.Sprintf("Getting content for deck w/ ID %s", deckID))

	var deckList *model.DeckList
	var err *cModel.APIError
	if deckList, err = skcDeckAPIDBInterface.GetDeckList(ctx, deckID); err != nil {
		err.HandleServerResponse(res)
		return
	}

	decodedListBytes, _ := base64.StdEncoding.DecodeString(deckList.ContentB64)
	decodedList := string(decodedListBytes) // decoded string of list contents

	var deckListBreakdown *model.DeckListBreakdown
	if deckListBreakdown, err = io.DeserializeDeckList(ctx, decodedList); err != nil {
		err.HandleServerResponse(res)
		return
	}
	deckList.MainDeck, deckList.ExtraDeck = deckListBreakdown.GetQuantities()

	logger.Info(fmt.Sprintf("Successfully retrieved deck list. Name {%s} and encoded deck list content {%s}. This deck list has {%d} main deck cards and {%d} extra deck cards.",
		deckList.Name, deckList.ContentB64, deckList.NumMainDeckCards, deckList.NumExtraDeckCards))
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(deckList)
}

func getDecksFeaturingCardHandler(res http.ResponseWriter, req *http.Request) {
	pathVars := mux.Vars(req)
	cardID := pathVars["cardID"]

	logger, ctx := cUtil.NewRequestSetup(context.Background(), "retrieve deck featuring card", slog.String("cardID", cardID))
	logger.Info("Fetching decks that feature card")

	suggestedDecks := model.SuggestedDecks{}

	suggestedDecks.FeaturedIn, _ = skcDeckAPIDBInterface.GetDecksThatFeatureCards(ctx, []string{cardID})

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(suggestedDecks)
}
