package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Info struct {
	ID   string `json:"id,omitempty"`
	From string `json:"from,omitempty"`
	Text string `json:"text,omitempty"`
}

// Slice for storing infos
var infos []Info

// Get info specified by id
func GetInfo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range infos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Info{})
}

// Get all infos
func GetInfos(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(infos)
}

// Append new info and get all infos
func PostInfo(w http.ResponseWriter, req *http.Request) {
	var info Info
	_ = json.NewDecoder(req.Body).Decode(&info)
	infos = append(infos, info)
	json.NewEncoder(w).Encode(infos)
}

// Delete info by id and get deleted info
func DeleteInfo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var deletedInfo Info
	for index, item := range infos {
		if item.ID == params["id"] {
			deletedInfo = Info{ID: item.ID, From: item.From, Text: item.Text}
			infos = append(infos[:index], infos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(deletedInfo)
}

func main() {
	// Create router to multiplex according to the path and the HTTP method
	router := mux.NewRouter()

	// Get all infos
	router.HandleFunc("/infos", GetInfos).Methods("GET")

	// Get info by id
	router.HandleFunc("/infos/{id}", GetInfo).Methods("GET")

	// Post info by id
	router.HandleFunc("/infos", PostInfo).Methods("POST")

	// Delete info by id
	router.HandleFunc("/infos/{id}", DeleteInfo).Methods("DELETE")

	// Listen and serve on localhost:22365
	log.Fatal(http.ListenAndServe(":22365", router))
}
