package goal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Start Goals Routes

func GoalRouteHandler(router *mux.Router) {
	goal := router.PathPrefix("/goal").Subrouter()
	goal.HandleFunc("", getGoals).Methods(http.MethodGet)
	goal.HandleFunc("", createGoal).Methods(http.MethodPost)
	goal.HandleFunc("/{goal:[0-9]+}", updateGoal).Methods(http.MethodPut)
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
		return
	}

	var resp, goalError = AddGoalToPersistence(newGoal)
	if goalError != nil {
		w.WriteHeader(http.StatusBadRequest)
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

func updateGoal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	goalId, queryErr := strconv.Atoi(vars["goal"])

	if queryErr != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := fmt.Sprintf(`{"error": "Method Not Found"}`)
		w.Write([]byte(errorResponse))
		return
	}

	decoder := json.NewDecoder(r.Body)

	var updatedGoal Goal
	bodyErr := decoder.Decode(&updatedGoal)

	if bodyErr != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := fmt.Sprintf(`{"error": "%s"}`, bodyErr.Error())
		w.Write([]byte(errorResponse))
		return
	}

	var resp, goalError = UpdateGoalInPersistence(updatedGoal, goalId)
	if goalError != nil {
		w.WriteHeader(http.StatusBadRequest)
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
