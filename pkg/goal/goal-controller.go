package goal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Start Goals Routes

func GoalRouteHandler(router *mux.Router) {
	goal := router.PathPrefix("/goal").Subrouter()
	goal.HandleFunc("", getGoals).Methods(http.MethodGet)
	goal.HandleFunc("", createGoal).Methods(http.MethodPost)
	goal.HandleFunc("", notFound)
}

func getGoals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	onlyActive := query.Get("onlyActive") == "true"

	var resp, goalError = GetGoalsFromPersistence(onlyActive)
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

func createGoal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var newGoal Goal
	bodyErr := decoder.Decode(&newGoal)

	if bodyErr != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := fmt.Sprintf(`{"error": "%s"}`, bodyErr.Error())
		w.Write([]byte(errorResponse))
	}

	var resp, goalError = AddGoalToPersistence(newGoal)
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
