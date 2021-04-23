package goal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Start Goals Routes

type Goal struct {
	Goal     int    `json:"goal"`
	GoalName string `json:"goalName"`
	Unit     string `json:"unit"`
	Active   bool   `json:"active"`
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

	var resp, goalError = GetGoalsFromPersistence()
	if goalError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := fmt.Sprintf(`{"error": "%s"}`, goalError.Error())
		w.Write([]byte(errorResponse))
		return
	}

	js, err := json.Marshal(resp)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := fmt.Sprintf(`{"error": "%s"}`, err.Error())
		w.Write([]byte(errorResponse))
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	if r.Method != "OPTIONS" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Method not found."}`))
	}
}

// End Goals Routes
