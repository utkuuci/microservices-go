package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

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

type RequestJSON struct {
	Action string      `json:"action"`
	Auth   AuthRequest `json:"auth,omitempty"`
	Log    LogRequest  `json:"log,omitempty"`
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

func (app *App) GetRequest(w http.ResponseWriter, r *http.Request) {
	var request RequestJSON
	err := ReadJson(r, &request)
	if err != nil {
		err = ErrorJson(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	switch request.Action {
	case "auth":
		data, _ := json.Marshal(request.Auth)
		req, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(data))
		if err != nil {
			err = ErrorJson(w, err)
			if err != nil {
				log.Println(err)
			}
			return
		}
		req.Header.Set("Content-Type", "application/json")
		r := &http.Client{}
		res, err := r.Do(req)

		if err != nil {
			err = ErrorJson(w, err)
			if err != nil {
				log.Println(err)
			}
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusAccepted {
			err = ErrorJson(w, err)
			if err != nil {
				log.Println(err)
			}
			return
		}

		var authRes Response
		err = json.NewDecoder(res.Body).Decode(&authRes)

		if err != nil {
			err = ErrorJson(w, err)
			if err != nil {
				log.Println(err)
			}
			return
		}

		response := Response{
			Error:   false,
			Message: "Authenticated user " + request.Auth.Email,
			Data:    authRes.Data,
		}
		err = WriteJson(w, http.StatusAccepted, response)
		if err != nil {
			log.Println(err.Error())
		}
	case "log":
		data, _ := json.Marshal(request.Log)

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
		var payload Response
		payload.Error = false
		payload.Message = "Logged via broker"
		WriteJson(w, http.StatusAccepted, payload)
	}
}

func (app *App) Broker(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Error:   false,
		Message: "Hit the broker",
	}

	err := WriteJson(w, http.StatusAccepted, res)
	if err != nil {
		log.Println(err)
	}
}
