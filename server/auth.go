package main

import (
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Auth is an object that holds incoming auth details
type Auth struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm-password"`
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	decoder := json.NewDecoder(r.Body)
	var auth Auth
	err := decoder.Decode(&auth)
	if err != nil {
		sendInternalError(w, err)
		return
	}

	username := auth.Username
	password := auth.Password
	confirmPassword := auth.ConfirmPassword

	if len(username) < 3 || len(password) < 3 {
		message := Message{http.StatusBadRequest, "Failure", "Username and Password must longer than 2 characters.", "", 0}
		sendMessageAndLogResult(w, message, nil)
		return
	}

	if password != confirmPassword {
		message := Message{http.StatusBadRequest, "Failure", "Passwords must be the same.", "", 0}
		sendMessageAndLogResult(w, message, nil)
		return
	}

	user, err := findUserInDb(username)
	if err != nil && err != mongo.ErrNoDocuments {
		sendInternalError(w, err)
		return
	}

	if (User{}) != user {
		message := Message{http.StatusBadRequest, "Failure", "Username already exists! Try a different username please.", "", 0}
		sendMessageAndLogResult(w, message, nil)
		return
	}

	hashedPassword := hashPassword(password)
	res, err := addUserToDb(username, hashedPassword)
	if err != nil || res == nil {
		sendInternalError(w, err)
		return
	}

	message := Message{http.StatusOK, "Success", "User has been added!", "", 0}
	sendMessageAndLogResult(w, message, nil)
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Print(err)
	}
	return string(bytes)
}

func sendInternalError(w http.ResponseWriter, err error) {
	message := Message{http.StatusInternalServerError, "Failure", "Internal issue, bug has been logged! Please try again later.", "", 0}
	sendMessageAndLogResult(w, message, err)
}

func sendMessageAndLogResult(w http.ResponseWriter, message Message, err error) {
	if err != nil {
		log.Print(err)
	} else {
		log.Print(message.Body)
	}
	sendMessage(w, message)
}

func sendMessage(w http.ResponseWriter, message Message) {
	w.WriteHeader(message.Status)
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Print(err)
	}
}
