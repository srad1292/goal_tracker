package progress

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Start Progress Routes

func ProgressRouteHandler(router *mux.Router) {
	progress := router.PathPrefix("/progress").Subrouter()
	progress.HandleFunc("/goal/{goalId:[0-9]+}", getProgress).Methods(http.MethodGet)
	progress.HandleFunc("", notFound)
}

func getProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	goalId, queryErr := strconv.Atoi(vars["goalId"])

	if queryErr != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := fmt.Sprintf(`{"error": "Method Not Found"}`)
		w.Write([]byte(errorResponse))
	}

	query := r.URL.Query()
	year, yearErr := strconv.Atoi(query.Get("year"))

	if yearErr != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := fmt.Sprintf(`{"error": "Query parameter, year, should be a number`)
		w.Write([]byte(errorResponse))
	}

	var resp, err = GetProgressFromPersistence(goalId, year)

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
		w.Write([]byte(`{"message": "not found"}`))
	}
}

// End Progress Routes
