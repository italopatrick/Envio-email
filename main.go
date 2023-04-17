package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Email struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Password string `json:"password"`
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into an Email struct
	var email Email
	err := json.NewDecoder(r.Body).Decode(&email)
	log.Println(email.Body)
	log.Println(email.Subject)
	log.Println(email.To)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return

	}

	// Set the SMTP server and port for your email provider
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"

	// Set up the authentication credentials
	auth := smtp.PlainAuth("", email.From, email.Password, smtpServer)

	// Set up the message body
	message := []byte("To: " + email.To + "\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"\r\n" +
		email.Body + "\r\n")

	// Send the email
	err = smtp.SendMail(smtpServer+":"+string(smtpPort), auth, email.From, []string{email.To}, message)
	if err != nil {
		http.Error(w, "Error sending email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Email sent successfully!")
}

func main() {
	// Set up the Gorilla mux router
	r := mux.NewRouter()
	credentials := handlers.AllowCredentials()
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"POST"})

	// start server listen
	// with error handling
	// Set up the send email endpoint
	r.HandleFunc("/send-email", sendEmailHandler).Methods("POST")

	// Start the server
	log.Println("Starting server on port 8000 ")
	err := http.ListenAndServe("8000", r)
	if err != nil {
		log.Fatal(http.ListenAndServe("PORT: 8000", handlers.CORS(origins, credentials, methods)(r)))
	}
}
