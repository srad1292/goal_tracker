package goal

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Start Goals Routes

type Goal struct {
	Goal     int    `json:"goal"`
	GoalName string `json:"goalName"`
	Unit     string `json:"unit"`
}

type GoalsResponse struct {
	Goals []Goal `json:"goals"`
}

func GoalRouteHandler(router *mux.Router) {
	goal := router.PathPrefix("/goal").Subrouter()
	goal.HandleFunc("", getGoals).Methods(http.MethodGet)
	goal.HandleFunc("", notFound)
}

func getGoals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp = GetGoalsFromPersistence()

	js, err := json.Marshal(resp)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "This is where I would return an error"}`))
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	if r.Method != "OPTIONS" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

// End Goals Routes
