// Core package used to configure skc-deck-api api and its endpoints.
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/ygo-skc/skc-deck-api/db"
	"github.com/ygo-skc/skc-deck-api/model"
	"github.com/ygo-skc/skc-deck-api/util"
)

const (
	CONTEXT = "/api/v1/deck"
)

var (
	skcDeckAPIDBInterface db.SKCDeckAPIDAO = db.SKCDeckAPIDAOImplementation{}
	router                *mux.Router
	corsOpts              *cors.Cors
	serverAPIKey          string
	chicagoLocation       *time.Location
)

func init() {
	// init Location
	if location, err := time.LoadLocation("America/Chicago"); err != nil {
		log.Fatalln("Could not load Chicago location")
	} else {
		chicagoLocation = location
	}
}

func verifyApiKey(headers http.Header) *model.APIError {
	clientKey := headers.Get("API-Key")

	if clientKey != serverAPIKey {
		log.Println("Client is using incorrect API Key. Cannot process request.")
		return &model.APIError{Message: "Request has incorrect or missing API Key."}
	}

	return nil
}

func verifyAPIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if err := verifyApiKey(req.Header); err != nil {
			res.Header().Add("Content-Type", "application/json")
			res.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(res).Encode(err)
		} else {
			next.ServeHTTP(res, req)
		}
	})
}

// sets common headers for response
func commonHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		res.Header().Add("Cache-Control", "max-age=300")

		next.ServeHTTP(res, req)
	})
}

// Configures routes and their middle wares
// This method should be called before the environment is set up as the API Key will be set according to the value found in environment
func ConfigureServer() {
	serverAPIKey = util.EnvMap["API_KEY"] // configure API Key

	router = mux.NewRouter()

	// configure non-admin routes
	unprotectedRoutes := router.PathPrefix(CONTEXT).Subrouter()
	unprotectedRoutes.HandleFunc("", submitNewDeckListHandler).Methods(http.MethodPost).Name("Deck List Submission")
	unprotectedRoutes.HandleFunc("/card/{cardID:[0-9]{8}}", getDecksFeaturingCardHandler).Methods(http.MethodGet).Name("Deck Featuring Card")
	unprotectedRoutes.HandleFunc("/{deckID:[0-9a-z]+}", getDeckListHandler).Methods(http.MethodGet).Name("Retrieve Info On Deck")

	// admin routes
	protectedRoutes := router.PathPrefix(CONTEXT).Subrouter()
	protectedRoutes.Use(verifyAPIKeyMiddleware)

	// common middleware
	router.Use(commonHeadersMiddleware)

	// Cors
	corsOpts = cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://dev.thesupremekingscastle.com", "https://dev.thesupremekingscastle.com", "https://thesupremekingscastle.com", "https://www.thesupremekingscastle.com"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodOptions,
		},

		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application
		},
	})
}

// configure server to handle HTTPS (secured) calls
func ServeTLS() {
	log.Println("Starting server in port 9001 (secured)")
	if err := http.ListenAndServeTLS(":9001", "certs/certificate.crt", "certs/private.key", corsOpts.Handler(router)); err != nil { // docker does not like localhost:9001 so im just using port number
		log.Fatalf("There was an error starting api server: %s", err)
	}
}

// configure server to handle HTTPs (un-secured) calls
func ServeUnsecured() {
	log.Println("Starting server in port 91 (unsecured)")
	if err := http.ListenAndServe(":91", corsOpts.Handler(router)); err != nil {
		log.Fatalf("There was an error starting api server: %s", err)
	}
}
