package state

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/AyuorusAguilar/chirpy/internal/database"
)

type State struct {
		Db *database.Queries
		Platform string
		fileserverHits atomic.Int32
		Secret string
		PKey string
	}

func (s *State) MdwareFileServerHitsBump(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.fileserverHits.Add(1)
		f.ServeHTTP(w, r)
	})
}

func (s *State) HandleGetFSH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	body := fmt.Appendf(nil, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, s.fileserverHits.Load())
	w.Write(body)
	// w.Write(fmt.Appendf(nil, "Hits: %d", s.fileserverHits.Load()))
}
func (s *State) HandleReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if s.Platform != "dev" {
		w.WriteHeader(403)
		w.Write(fmt.Appendf(nil, "403 Forbidden"))
		return
	}
	s.fileserverHits.Store(0)
	s.Db.Reset(context.Background())
	w.WriteHeader(200)
	w.Write(fmt.Appendf(nil, "	Metrics and database have been reset"))
}

func NewState(db *database.Queries) State {
	return State{
		Db: db,
		fileserverHits: atomic.Int32{},
	}
}