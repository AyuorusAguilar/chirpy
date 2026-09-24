package state

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type State struct {
		fileserverHits atomic.Int32
	}

func (s *State) MdwareFileServerHitsBump(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.fileserverHits.Add(1)
		f.ServeHTTP(w, r)
	})
}

func (s *State) HandleGetFSH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(fmt.Appendf(nil, "Hits: %d", s.fileserverHits.Load()))
}
func (s *State) HandleReset(w http.ResponseWriter, r *http.Request) {
	s.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(fmt.Appendf(nil, "Metrics have been reset"))
}

func NewState() State {
	return State{
		fileserverHits: atomic.Int32{},
	}
}