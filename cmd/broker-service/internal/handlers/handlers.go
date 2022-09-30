package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action`
	Auth   AuthPayload `json:"auth, omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:password`
}

func Broker(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = WriteJSON(w, http.StatusOK, payload)
}

func HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := ReadJSON(w, r, &requestPayload)
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		// app.authenticate(w, requestPayload.Auth)
	default:
		ErrorJSON(w, errors.New("unknown action"))
	}
}

func authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService JsonResponse

	//decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	WriteJSON(w, http.StatusAccepted, payload)
}
