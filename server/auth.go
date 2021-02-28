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

	validInputs := checkInputs(w, auth)
	if !validInputs {
		return
	}

	hashedPassword := hashPassword(w, auth)
	if hashedPassword == "" {
		return
	}

	res, err := addUserToDb(auth.Username, hashedPassword)
	if err != nil || res == nil {
		sendInternalError(w, err)
		return
	}

	message := Message{http.StatusOK, "Success", "User has been added!", "", 0}
	sendMessageAndLogResult(w, message, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	decoder := json.NewDecoder(r.Body)
	var auth Auth
	err := decoder.Decode(&auth)
	if err != nil {
		sendInternalError(w, err)
		return
	}

	user, err := findUserInDb(auth.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		sendInternalError(w, err)
		return
	}

	if (User{}) == user {
		message := Message{http.StatusBadRequest, "Failure", "Username " + auth.Username + " does not exist!", "", 0}
		sendMessageAndLogResult(w, message, nil)
		return
	}

	validPassword := checkPassword(w, auth.Password, user.Password)
	if !validPassword {
		return
	}

	message := Message{http.StatusOK, "Success", "User " + user.Username + " has logged in!", "", 0}
	sendMessageAndLogResult(w, message, nil)
}

func checkInputs(w http.ResponseWriter, auth Auth) bool {
	username := auth.Username
	password := auth.Password
	confirmPassword := auth.ConfirmPassword

	if len(username) < 3 || len(password) < 3 {
		message := Message{http.StatusBadRequest, "Failure", "Username and Password must longer than 2 characters.", "", 0}
		sendMessageAndLogResult(w, message, nil)
		return false
	}

	if password != confirmPassword {
		message := Message{http.StatusBadRequest, "Failure", "Passwords must be the same.", "", 0}
		sendMessageAndLogResult(w, message, nil)
		return false
	}

	user, err := findUserInDb(username)
	if err != nil && err != mongo.ErrNoDocuments {
		sendInternalError(w, err)
		return false
	}

	if (User{}) != user {
		message := Message{http.StatusBadRequest, "Failure", "Username " + user.Username + " already exists! Try a different username please.", "", 0}
		sendMessageAndLogResult(w, message, nil)
		return false
	}

	return true
}

func hashPassword(w http.ResponseWriter, auth Auth) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(auth.Password), 14)
	if err != nil {
		sendInternalError(w, err)
		return ""
	}
	return string(bytes)
}

func checkPassword(w http.ResponseWriter, attemptedPassword string, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(attemptedPassword))
	if err != nil {
		sendInternalError(w, err)
		return false
	}
	return true
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
