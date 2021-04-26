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
	progress.HandleFunc("/time/goal/{goalId:[0-9]+}", getProgressByPeriod).Methods(http.MethodGet)
	progress.HandleFunc("/goal/{goalId:[0-9]+}/session", getBestSessions).Methods(http.MethodGet)
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

func getProgressByPeriod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	goalId, queryErr := strconv.Atoi(vars["goalId"])

	if queryErr != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := fmt.Sprintf(`{"error": "Method Not Found"}`)
		w.Write([]byte(errorResponse))
		return
	}

	query := r.URL.Query()
	year, yearErr := strconv.Atoi(query.Get("year"))

	if yearErr != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := fmt.Sprintf(`{"error": "Query parameter, year, should be a number`)
		w.Write([]byte(errorResponse))
		return
	}

	period := query.Get("period")

	if period == "" || (period != "year" && period != "month" && period != "day") {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := fmt.Sprintf(`{"error": "Query parameter 'year' should be one of: year, month, day`)
		w.Write([]byte(errorResponse))
		return
	}

	edge := query.Get("edge")

	if edge != "" && !(edge == "high" || edge == "low") {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := fmt.Sprintf(`{"error": "Query parameter 'edge', if used, should be one of: high, low`)
		w.Write([]byte(errorResponse))
		return
	}

	var resp ProgressByTimeResponse
	var err error
	if edge == "" {
		resp, err = GetProgressByTimeFromPersistence(goalId, year, period)
	} else {
		resp, err = GetBestProgressByTimeFromPersistence(goalId, year, period, edge)
	}

	js, err := json.Marshal(resp)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := fmt.Sprintf(`{"error": "%s"}`, err.Error())
		w.Write([]byte(errorResponse))
		return
	}
}

func getBestSessions(w http.ResponseWriter, r *http.Request) {
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

	useLow := query.Get("useLow") == "true"

	var resp, err = GetBestSessionsFromPersistence(goalId, year, useLow)

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
