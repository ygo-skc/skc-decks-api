package db

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ygo-skc/skc-deck-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	skcDeckDB          *mongo.Database
	deckListCollection *mongo.Collection
)

// interface
type SKCDeckAPIDAO interface {
	GetSKCDeckAPIDBVersion() (string, error)

	InsertDeckList(deckList model.DeckList)
	GetDeckList(deckID string) (*model.DeckList, *model.APIError)
	GetDecksThatFeatureCards([]string) (*[]model.DeckList, *model.APIError)
}

// impl
type SKCDeckAPIDAOImplementation struct{}

// Retrieves the version number of the SKC Deck API DB or throws an error if an exception occurs.
func (dbInterface SKCDeckAPIDAOImplementation) GetSKCDeckAPIDBVersion() (string, error) {
	var commandResult bson.M
	command := bson.D{{Key: "serverStatus", Value: 1}}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := skcDeckDB.RunCommand(ctx, command).Decode(&commandResult); err != nil {
		log.Println("Error getting SKC Deck API DB version", err)
		return "", err
	} else {
		return fmt.Sprintf("%v", commandResult["version"]), nil
	}
}

func (dbInterface SKCDeckAPIDAOImplementation) InsertDeckList(deckList model.DeckList) {
	deckList.CreatedAt = time.Now()
	deckList.UpdatedAt = deckList.CreatedAt

	log.Printf("Inserting deck with name %s with Main Deck size %d and Extra Deck size %d. List contents (in base64 and possibly reformatted) %s",
		deckList.Name, deckList.NumMainDeckCards, deckList.NumExtraDeckCards, deckList.ContentB64)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if res, err := deckListCollection.InsertOne(ctx, deckList); err != nil {
		log.Println("Error inserting new deck list into DB", err)
	} else {
		log.Println("Successfully inserted new deck list into DB, ID:", res.InsertedID)
	}
}

func (dbInterface SKCDeckAPIDAOImplementation) GetDeckList(deckID string) (*model.DeckList, *model.APIError) {
	if objectId, err := primitive.ObjectIDFromHex(deckID); err != nil {
		log.Println("Invalid Object ID.")
		return nil, &model.APIError{Message: "Object ID used for deck list was not valid.", StatusCode: http.StatusBadRequest}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		var dl model.DeckList
		if err := deckListCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&dl); err != nil {
			log.Printf("Error retrieving deck list w/ ID %s. Err: %v", deckID, err)
			return nil, &model.APIError{Message: "Error retrieving deck", StatusCode: http.StatusNotFound}
		} else {
			return &dl, nil
		}
	}
}

func (dbInterface SKCDeckAPIDAOImplementation) GetDecksThatFeatureCards(cardIDs []string) (*[]model.DeckList, *model.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// select only these fields from collection
	opts := options.Find().SetProjection(
		bson.D{
			{Key: "name", Value: 1}, {Key: "videoUrl", Value: 1}, {Key: "uniqueCards", Value: 1}, {Key: "deckMascots", Value: 1}, {Key: "numMainDeckCards", Value: 1},
			{Key: "numExtraDeckCards", Value: 1}, {Key: "tags", Value: 1}, {Key: "createdAt", Value: 1}, {Key: "updatedAt", Value: 1},
		},
	)

	if cursor, err := deckListCollection.Find(ctx, bson.M{"uniqueCards": bson.M{"$in": cardIDs}}, opts); err != nil {
		log.Printf("Error retrieving all deck lists that feature cards w/ ID %v. Err: %v", cardIDs, err)
		return nil, &model.APIError{Message: "Error retrieving deck suggestions", StatusCode: http.StatusInternalServerError}
	} else {
		dl := []model.DeckList{}
		if err := cursor.All(ctx, &dl); err != nil {
			log.Printf("Error retrieving all deck lists that feature cards w/ ID %v. Err: %v", cardIDs, err)
			return nil, &model.APIError{Message: "Error retrieving deck suggestions", StatusCode: http.StatusInternalServerError}
		}

		return &dl, nil
	}
}
