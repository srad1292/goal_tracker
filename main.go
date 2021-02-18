package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Middleware
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "OPTIONS" {
			log.Println(request.Method, request.URL)
			next.ServeHTTP(writer, request)
		}
	})
}

// Base Route
func base(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Thanks for hitting the goal tracker API!"}`))
}

// Start Progress Routes
func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post called"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "put called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	if r.Method != "OPTIONS" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

// End Progress Routes

// Start Goals Routes

type goal struct {
	Goal     int    `json:"goal"`
	GoalName string `json:"goalName"`
	Unit     string `json:"unit"`
}

type goalsResponse struct {
	Goals []goal `json:"goals"`
}

func getGoals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp = goalsResponse{
		Goals: []goal{
			{
				Goal:     1,
				GoalName: "Push Ups",
				Unit:     "",
			},
			{
				Goal:     2,
				GoalName: "Drawing",
				Unit:     "Minutes",
			},
		},
	}

	js, err := json.Marshal(resp)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "This is where I would return an error"}`))
	}
}

// End Goals Routes

func main() {
	router := mux.NewRouter()
	router.Use(logRequest)

	router.HandleFunc("/", base).Methods(http.MethodGet)
	router.HandleFunc("/", notFound)

	progress := router.PathPrefix("/progress").Subrouter()
	progress.HandleFunc("", get).Methods(http.MethodGet)
	progress.HandleFunc("", post).Methods(http.MethodPost)
	progress.HandleFunc("", put).Methods(http.MethodPut)
	progress.HandleFunc("", delete).Methods(http.MethodDelete)
	progress.HandleFunc("", notFound)

	goal := router.PathPrefix("/goal").Subrouter()
	goal.HandleFunc("", getGoals).Methods(http.MethodGet)
	goal.HandleFunc("", notFound)

	log.Println("Starting server at localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
