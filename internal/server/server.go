package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID       int
	Username string
}

type Server struct {
	sync.RWMutex
	db     sync.Map
	cache  sync.Map
	dbhits int
}

func NewServer() *Server {
	var db sync.Map
	var cache sync.Map
	for i := 0; i < 100; i++ {
		db.Store(i+1, &User{
			ID:       i + 1,
			Username: fmt.Sprintf("user_%d", i),
		})
	}
	return &Server{
		db:    db,
		cache: cache,
	}
}

func (s *Server) tryCache(id int) (any, bool) {
	s.RLock()
	defer s.RUnlock()
	user, ok := s.cache.Load(id)
	if ok {
		return user, true
	}
	return nil, false
}

func (s *Server) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	if idstr == "" || len(idstr) == 0 {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, "id conversion failed", http.StatusBadRequest)
		return
	}

	// try the cache
	user, ok := s.tryCache(id)
	if ok {
		json.NewEncoder(w).Encode(user)
		return
	}

	//hit the database
	user, ok = s.db.Load(id)
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	s.dbhits++

	// insert the cache
	s.cache.Store(id, user)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "json encode failed", http.StatusInternalServerError)
		return
	}
}
