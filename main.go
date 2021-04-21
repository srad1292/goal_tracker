package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/srad1292/goal_tracker/pkg/database"
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
	readEnvironment()
	validateEnvironment()

	database.GetDatabase()

	router := mux.NewRouter()
	router.Use(logRequest)

	router.HandleFunc("/", base).Methods(http.MethodGet)
	router.HandleFunc("/", notFound)

	goal.GoalRouteHandler(router)
	progress.ProgressRouteHandler(router)

	log.Println("Starting server at localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}

func readEnvironment() {
	error := godotenv.Load(".env")

	if error != nil {
		log.Println(error)
		log.Fatalf("Failed to load .env file")
	}

}

func validateEnvironment() {
	missing := make([]string, 0)

	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	db_schema := os.Getenv("DB_SCHEMA")
	db_conn := os.Getenv("DB_MAX_CONNECTIONS")

	if db_user == "" {
		missing = append(missing, "DB_USER")
	}

	if db_password == "" {
		missing = append(missing, "DB_PASSWORD")
	}

	if db_host == "" {
		missing = append(missing, "DB_HOST")
	}

	if db_port == "" {
		missing = append(missing, "DB_PORT")
	}

	if db_name == "" {
		missing = append(missing, "DB_NAME")
	}

	if db_schema == "" {
		missing = append(missing, "DB_SCHEMA")
	}

	if db_conn == "" {
		missing = append(missing, "DB_MAX_CONNECTIONS")
	}

	if len(missing) > 0 {
		log.Println("MISSING ENVIRONMENT VARIABLES:")
		for _, value := range missing {
			log.Printf("%s, ", value)
		}
		log.Println()
	}

}
