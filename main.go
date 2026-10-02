package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/AyuorusAguilar/chirpy/internal/api"
	"github.com/AyuorusAguilar/chirpy/internal/state"
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
	s.Secret = os.Getenv("SICRET")
	s.PKey = os.Getenv("POLKAKEY")
	mux := http.NewServeMux()
	serv := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	/* Routing */
	mux.Handle("GET /app/", http.StripPrefix("/app/", s.MdwareFileServerHitsBump(http.FileServer(http.Dir(".")))))
	
	mux.Handle("GET /admin/metrics", http.HandlerFunc(s.HandleGetFSH))
	mux.Handle("POST /admin/reset", http.HandlerFunc(s.HandleReset))

	mux.HandleFunc("GET  	/api/healthz", api.Ok)
	mux.HandleFunc("POST 	/api/chirps", api.MdwareAccessState(&s, api.PostChirps))
	mux.HandleFunc("GET  	/api/chirps/", api.MdwareAccessState(&s, api.GetChirps))
	mux.HandleFunc("GET 	/api/chirps/{ChirpId}", api.MdwareAccessState(&s, api.GetChirpsById))
	mux.HandleFunc("DELETE 	/api/chirps/{ChirpId}", api.MdwareAccessState(&s, api.DeleteChirp))

	mux.HandleFunc("POST 	/api/users", api.MdwareAccessState(&s, api.PostUser))
	mux.HandleFunc("PUT	 	/api/users", api.MdwareAccessState(&s, api.PutUser))
	mux.HandleFunc("POST 	/api/login", api.MdwareAccessState(&s, api.PostLogin))
	mux.HandleFunc("POST 	/api/refresh", api.MdwareAccessState(&s, api.RefreshToken))
	mux.HandleFunc("POST 	/api/revoke", api.MdwareAccessState(&s, api.RevokeToken))
	
	mux.HandleFunc("POST 	/api/polka/webhooks", api.MdwareAccessState(&s, api.UpgradeUser))
	
	/* Start */
	fmt.Printf("Opening server to listen at port 8080\n")
	serv.ListenAndServe()
}
