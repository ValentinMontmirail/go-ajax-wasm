package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"wasm/pkg/authors"
	"wasm/pkg/database"
	"wasm/pkg/middleware"
	"wasm/pkg/token"

	log "github.com/sirupsen/logrus"
)

func getAllAuthors(db *database.MyDB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Trace("getAllAuthors -- IN")

	author, err := db.GetAllAuthors()
	if err != nil {
		fmt.Fprintf(w, "GetAllAuthors() err: %v", err)
		return
	}

	if len(author) == 0 {
		fmt.Fprintf(w, "[]")
		return
	}

	authorJSON, err := json.Marshal(author)
	if err != nil {
		fmt.Fprintf(w, "Marshal() err: %v", err)
		return
	}

	fmt.Fprintf(w, string(authorJSON))
	log.Trace("getAllAuthors -- OUT")
}

func createAuthor(db *database.MyDB, w http.ResponseWriter, r *http.Request) {

	log.Trace("createAuthor -- IN")

	// Decode the JSON body into the struct
	var authorReq authors.CreateAuthorParams
	err := json.NewDecoder(r.Body).Decode(&authorReq)
	if err != nil {
		log.Errorf("Error decoding request body: %v", err)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	author, err := db.CreateAuthor(authorReq.Name, authorReq.Bio)
	if err != nil {
		log.Errorf("Error decoding request body: %v", err)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	authorJSON, err := json.Marshal(author)
	if err != nil {
		log.Errorf("Marshal() err: %v", err)
		fmt.Fprintf(w, "Marshal() err: %v", err)
		return
	}

	log.Trace("createAuthor -- OUT")

	fmt.Fprintf(w, string(authorJSON))
}

func init() {
	// Set the logger's formatter to logrus.TextFormatter with specific settings
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,         // Force colored output
		FullTimestamp:   true,         // Enable full timestamp in logs
		TimestampFormat: time.RFC3339, // Set the format of the timestamp
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	ctx := context.Background()

	http.Handle("/", middleware.ApplyMiddleware(fs)) // Apply CSP, the rate limiter and the CORS to the file server

	db, err := database.Open(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.Handle("/api/v1/token", middleware.ApplyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Trace("getToken() -- IN")

		currentToken, expDate := token.GetCurrentToken()
		// Set the token in a secure, HTTP-only cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "AuthToken",
			Value:    currentToken,
			Expires:  expDate,                 // Expires in the next 10 seconds
			Secure:   true,                    // Only send over HTTPS
			HttpOnly: true,                    // Not accessible via JavaScript
			Path:     "/api/v1/",              // Accessible on APIs only
			SameSite: http.SameSiteStrictMode, // Strict same-site policy
		})

		log.Trace("getToken() -- OUT")
	})))

	http.Handle("/api/v1/authors", middleware.ApplyTokenValidator(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Add("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			getAllAuthors(db, w, r)
		case "POST":
			createAuthor(db, w, r)
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	})))

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cmd/webserver/server-cert.pem", "cmd/webserver/server-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Create a Server instance with TLS Config
	server := &http.Server{
		Addr: ":3000",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{serverCert},
		},
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on https://localhost:3000/index.html")

	// Start the server with TLS
	log.Fatal(server.ListenAndServeTLS("", ""))
}
