package main

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load env vars
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GarrettKeith Email Relay Server")
	})

	http.HandleFunc("/email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.FormValue("to") != "" && r.FormValue("subject") != "" && r.FormValue("body") != "" {
			// Check if key and secret are valid
			key := r.FormValue("key")
			secret := r.FormValue("secret")
			if !authenticate(key, secret) {
				// Return authentication error
				fmt.Fprintf(w, "Error authenticating")
				return
			}

			fromemail := ""
			fromname := r.FormValue("fromname")
			if fromname != "" {
				fromemail = fromname + " <" + fromemail + ">"
			}
			password := ""

			// Get values from post form
			to := r.FormValue("to")
			subject := r.FormValue("subject")
			body := r.FormValue("body")

			// Build message
			msg := "From: " + fromemail + "\n" + "To: " + to + "\n" + "Subject: " + subject + "\n\n" + body + "\n"

			// Send email
			err := smtp.SendMail("smtp.zoho.com:587", smtp.PlainAuth("", fromemail, password, "smtp.zoho.com"), fromemail, []string{to}, []byte(msg))
			if err != nil {
				fmt.Println(err)
				// Return error
				fmt.Fprintf(w, "Error sending email")
				return
			}

			fmt.Println("Email Sent!")
			fmt.Fprintf(w, "Email Sent!")
			return
		}

		// Return error
		fmt.Fprintf(w, "Error sending email, missing required fields. Please include to, subject, and body.")
	})

	http.ListenAndServe(":8080", nil)
}

func authenticate(key string, secret string) bool {
	// Check if key and secret are valid from the .env file
	if key == os.Getenv("API_KEY") && secret == os.Getenv("API_SECRET") {
		return true
	}
	return false
}
