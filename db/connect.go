package db

import (
	"fmt"
	"log"
	"time"

	"github.com/ygo-skc/skc-deck-api/util"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

const (
	minPoolSize = 20
	maxPoolSize = 30
)

// connects to Deck API database
func EstablishSKCDeckAPIDBConn() {
	certificateKeyFilePath := "./certs/skc-deck-api-db.pem"
	uri := fmt.Sprintf("%s/?tlsCertificateKeyFile=%s", util.EnvMap["DB_HOST"], certificateKeyFilePath)

	credential := options.Credential{
		AuthMechanism: "MONGODB-X509",
	}

	if client, err := mongo.Connect(options.Client().
		ApplyURI(uri).
		SetAuth(credential).
		SetMaxPoolSize(maxPoolSize).
		SetMaxConnIdleTime(10 * time.Minute).
		SetTimeout(2 * time.Second).
		SetReadConcern(readconcern.Majority()).   // prefer strongly consistent reeds
		SetWriteConcern(writeconcern.Majority()). // writes to most replicas before acknowledging the write is complete
		SetAppName("SKC Deck API")); err != nil {
		log.Fatalln("Error creating new mongodb client for skc-deck-api-db", err)
	} else {
		skcDeckDB = client.Database("deckDB")
	}

	// init collections
	deckListCollection = skcDeckDB.Collection("lists")
}
