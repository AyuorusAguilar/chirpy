package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AyuorusAguilar/chirpy/Internal/api"
	"github.com/AyuorusAguilar/chirpy/Internal/state"
	"github.com/AyuorusAguilar/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
	
func main() {
	/* Initialization */
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error starting db connection: %v", err)
		os.Exit(1)
	}
	s := state.NewState(database.New(db))
	s.Platform = os.Getenv("PLATFORM")
	mux := http.NewServeMux()
	serv := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	/* Routing */
	mux.Handle("GET /app/", http.StripPrefix("/app/", s.MdwareFileServerHitsBump(http.FileServer(http.Dir(".")))))
	
	mux.Handle("GET /admin/metrics", http.HandlerFunc(s.HandleGetFSH))
	mux.Handle("POST /admin/reset", http.HandlerFunc(s.HandleReset))

	mux.HandleFunc("GET /api/healthz", api.Ok)
	mux.HandleFunc("POST /api/chirps", api.MdwareAccessState(&s, api.PostChirps))
	mux.HandleFunc("GET /api/chirps", api.MdwareAccessState(&s, api.GetChirps))
	mux.HandleFunc("GET /api/chirps/{GetChirpsById}", api.MdwareAccessState(&s, api.GetChirpsById))
	mux.HandleFunc("POST /api/users", api.MdwareAccessState(&s, api.CreateUser))
	

	/* Start */
	fmt.Printf("Opening server to listen at port 8080\n")
	serv.ListenAndServe()
}
