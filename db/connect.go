package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	minPoolSize = 20
	maxPoolSize = 40
)

func EstablishSKCSuggestionEngineDBConn() {
	certificateKeyFilePath := "./certs/skc-suggestion-engine-db.pem"
	uri := "mongodb+srv://skc-suggestion-engine-e.rfait.mongodb.net/?tlsCertificateKeyFile=%s"
	uri = fmt.Sprintf(uri, certificateKeyFilePath)

	credential := options.Credential{
		AuthMechanism: "MONGODB-X509",
	}

	if client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetAuth(credential).SetMinPoolSize(minPoolSize).SetMaxPoolSize(maxPoolSize).SetMaxConnIdleTime(10*time.Minute).
		SetAppName("SKC Suggestion Engine")); err != nil {
		log.Fatalln("Error creating new mongodb client for skc-suggestion-engine DB", err)
	} else {
		skcSuggestionDB = client.Database("suggestionDB")
	}

	// init collections
	deckListCollection = skcSuggestionDB.Collection("deckLists")
}
