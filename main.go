package main

import (
	"log"
	"net/http"

	"time"

	"github.com/gorilla/mux"
	"toonai.go/handlers"
	"toonai.go/middleware"
)
func main() {
	r := mux.NewRouter()

	// Initialize routes
	initializeRoutes(r)

	// Start server
	log.Printf("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initializeRoutes(r *mux.Router) {
	// Create a subrouter for /api
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/stream")
		for i := 0; i < 10; i++ {
			w.Write([]byte("Hello, World!\n"))
			w.(http.Flusher).Flush()
			time.Sleep(1 * time.Second)
		}
		// w.Write([]byte("Hello, World!\n"))
	}).Methods("GET")

	api.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello, World!"}`))
	})

	// Comic routes
	api.HandleFunc("/comics/chapters/{chapterID}", handlers.GetComicChapter).Methods("GET")
	
	// Video routes
	api.HandleFunc("/videos/{videoID}", handlers.StreamVideo).Methods("GET")

	// Add middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
} 