package sessions

import (
	"code.google.com/p/gorilla/sessions"
	"net/http"
)

const(
  _ = iota
  Anonymous
  User
)

var (
	store = sessions.NewCookieStore([]byte("NiseGoPostSecret"))
)

type Session struct {
	*sessions.Session
}

func Get(r *http.Request) *Session {
	s, _ := store.Get(r, "session")
	return &Session{s}
}

func New(r *http.Request) *Session {
	s, _ := store.New(r, "session")
	s.Values["hasError"] = false
	return &Session{s}
}

func (s Session) ClearFlashes() {
	s.Values["_flash"] = make([]interface{}, 0)
}
