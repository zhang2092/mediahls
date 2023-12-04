package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Respond(w http.ResponseWriter, message string, v any, statusCode int) {
	rsp := response{
		Success: true,
		Message: message,
		Data:    v,
	}
	b, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(b)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("could not write http response: %v\n", err)
	}
}

func RespondErr(w http.ResponseWriter, message string, v any) {
	rsp := response{
		Success: false,
		Message: message,
		Data:    v,
	}
	b, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(b)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("could not write http response: %v\n", err)
	}
}

func RespondJson(w http.ResponseWriter, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("could not write http response: %v\n", err)
	}
}
