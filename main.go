package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/srad1292/goal_tracker/pkg/goal"
	"github.com/srad1292/goal_tracker/pkg/progress"
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

func main() {
	router := mux.NewRouter()
	router.Use(logRequest)

	router.HandleFunc("/", base).Methods(http.MethodGet)
	router.HandleFunc("/", notFound)

	goal.GoalRouteHandler(router)
	progress.ProgressRouteHandler(router)

	log.Println("Starting server at localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
