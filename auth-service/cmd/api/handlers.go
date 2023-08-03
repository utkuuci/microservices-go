package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"log"
	"net/http"
	"time"
)

var SECRET = []byte("super-secret")

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LogRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func ReadJson(r *http.Request, data interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must be have only a single JSON value")
	}
	return nil
}

func WriteJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func ErrorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload Response
	payload.Error = true
	payload.Message = err.Error()

	return WriteJson(w, statusCode, payload)
}

func createJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenStr, err := token.SignedString(SECRET)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (app *App) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	err := ReadJson(r, &req)
	if err != nil {
		err = ErrorJson(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	user, err := app.Models.AuthenticationModel.GetByEmail(req.Email)
	log.Println(user)
	if err != nil {
		err = ErrorJson(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if user.Password != req.Password {
		err = ErrorJson(w, errors.New("Invalid email or password"))
		if err != nil {
			log.Println(err)
		}
		return
	}
	tok, err := createJWT(req.Email)
	if err != nil {
		ErrorJson(w, err)
		if err != nil {
			log.Println(err)
		}
	}
	LoggingRequest(w, LogRequest{
		Name: "authentication",
		Data: time.Now().String() + " logged in user email: " + req.Email,
	})
	var response Response
	response.Error = false
	response.Message = "User authenticated " + req.Email
	response.Data = struct {
		Token string `json:"token"`
	}{
		tok,
	}
	err = WriteJson(w, http.StatusAccepted, response)
	if err != nil {
		log.Println(err)
	}
}

func LoggingRequest(w http.ResponseWriter, entry LogRequest) {
	data, _ := json.Marshal(entry)

	req, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(data))
	if err != nil {
		err = ErrorJson(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		err = ErrorJson(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		if err != nil {
			err = ErrorJson(w, err)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
}
