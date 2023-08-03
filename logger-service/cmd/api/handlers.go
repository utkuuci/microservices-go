package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	data2 "logger-serivce/data"
	"net/http"
)

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RequestPayload struct {
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
	log.Println(out)
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

func (app *App) Log(w http.ResponseWriter, r *http.Request) {
	var request RequestPayload
	err := ReadJson(r, &request)

	if err != nil {
		err = ErrorJson(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	event := data2.Log{
		Name: request.Name,
		Data: request.Data,
	}

	err = app.Models.Log.Insert(event)
	if err != nil {
		err = ErrorJson(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	response := Response{
		Error:   false,
		Message: "Logged",
	}

	err = WriteJson(w, http.StatusAccepted, response)
	if err != nil {
		err = ErrorJson(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
}
