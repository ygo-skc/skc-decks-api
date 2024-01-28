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

// connects to Deck API database
func EstablishSKCDeckAPIDBConn() {
	certificateKeyFilePath := "./certs/skc-deck-api-db.pem"
	uri := "mongodb+srv://skc-deck-api-east1.twjikpk.mongodb.net/?tlsCertificateKeyFile=%s"
	uri = fmt.Sprintf(uri, certificateKeyFilePath)

	credential := options.Credential{
		AuthMechanism: "MONGODB-X509",
	}

	if client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri).
		SetAuth(credential).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize).
		SetMaxConnIdleTime(10*time.Minute).
		SetAppName("SKC Deck API")); err != nil {
		log.Fatalln("Error creating new mongodb client for skc-deck-api-db", err)
	} else {
		skcDeckDB = client.Database("deckDB")
	}

	// init collections
	deckListCollection = skcDeckDB.Collection("lists")
}
