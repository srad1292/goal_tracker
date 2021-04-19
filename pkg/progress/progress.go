package progress

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Start Progress Routes

type ProgressPersistence struct {
	Progress    int    `json:"progress"`
	Amount      int    `json:"amount"`
	SessionDate string `json:"sessionDate"`
	Goal        int    `json:"goal"`
	GoalName    string `json:"goalName"`
	Unit        string `json:"unit"`
}

type ProgressResponse struct {
	Progress []ProgressPersistence `json:"progress"`
}

func ProgressRouteHandler(router *mux.Router) {
	progress := router.PathPrefix("/progress").Subrouter()
	progress.HandleFunc("", getProgress).Methods(http.MethodGet)
	progress.HandleFunc("", notFound)
}

func getProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp = GetProgressFromPersistence()

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

// End Progress Routes
