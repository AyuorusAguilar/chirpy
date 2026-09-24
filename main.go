package main

import (
	"fmt"
	"net/http"
	"github.com/AyuorusAguilar/chirpy/Internal/state"
)
	
func main() {
	s := state.NewState()
	mux := http.NewServeMux()
	serv := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	mux.Handle("/app/", http.StripPrefix("/app/", s.MdwareFileServerHitsBump(http.FileServer(http.Dir(".")))))
	mux.Handle("/healthz", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(200)
			w.Write([]byte("OK"))	
		}))
	mux.Handle("/metrics", http.HandlerFunc(s.HandleGetFSH))
	mux.Handle("/reset", http.HandlerFunc(s.HandleReset))
	fmt.Printf("Opening server to listen at port 8080\n")
	serv.ListenAndServe()
}
