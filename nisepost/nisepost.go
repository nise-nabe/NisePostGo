package nisepost

import (
	"code.google.com/p/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
)

type NisePostGoHandler struct {
	handler func(http.ResponseWriter, *http.Request, *sessions.Session)
}

func (h *NisePostGoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	s, _ := store.Get(r, "session")
	h.handler(w, r, s)
}

func Init() {
	store = sessions.NewCookieStore([]byte("NiseGoPostSecret"))
	initDB()
	initRouting()
}

var (
	store *sessions.CookieStore
)

func initRouting() {
	http.Handle("/", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		if r.URL.Path != "/" {
			t := LoadWebContent(r.URL.Path)
			t.Execute(w, nil)
			return
		}
		log.Println(s.Values["name"])
		log.Println(s.Values["role"])
		t := LoadTemplate("index.html")
		t.Execute(w, nil)
	}})
	http.Handle("/login", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		if s.Values["role"] == nil {
			s, _ = store.New(r, "session")
			s.Values["role"] = "Anonymous"
		} else if s.Values["role"] != "Anonymous" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		s.Save(r, w)
		t := LoadTemplate("login.html")
		t.Execute(w, s)
	}})
	http.Handle("/login/post", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		switch r.Method {
		case "POST":
			if s.Values["role"] != "Anonymous" {
				break
			}
			username, password := r.FormValue("username"), r.FormValue("password")
			if !Authenticate(username, password) {
				break
			}
			s.Values["name"] = username
			s.Values["role"] = "User"
			s.Save(r, w)
			log.Println("User Authorized")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		s.Values["_flash"] = make([]interface{}, 0)
		s.AddFlash("login was not succeeded!")
		s.Save(r, w)
		log.Println("User Unauthorized")
		http.Redirect(w, r, "/login", http.StatusFound)
	}})
	http.Handle("/logout", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		s, _ = store.New(r, "session")
		s.Values["role"] = "Anonymous"
		s.Values["_flash"] = make([]interface{}, 0)
		s.AddFlash("logout was succeeded!")
		s.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusFound)
	}})
	http.Handle("/register", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		switch r.Method {
		case "GET":
			if s.Values["role"] != "Anonymous" {
				break
			}
			t := LoadTemplate("register.html")
			t.Execute(w, s)
		}
	}})
	http.Handle("/register/post", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		if s.Values["role"] != "Anonymous" {
			return
		}
		switch r.Method {
		case "POST":
			username := r.FormValue("username")
			password, password2 := r.FormValue("password"), r.FormValue("password2")
			if password == "" {
				s.AddFlash("password was empty")
				http.Redirect(w, r, "/register", http.StatusFound)
			}
			if password2 == "" {
				s.AddFlash("confirmed password was empty")
				http.Redirect(w, r, "/register", http.StatusFound)
			}
			if password != password2 {
				s.AddFlash("password doesn't correspond")
				http.Redirect(w, r, "/register", http.StatusFound)
			}

			NewUser(username, password).Save()
			s.AddFlash("Register was succeeded. ")
			s.Values["name"] = username
			s.Values["role"] = "User"
			s.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}})
}

func LoadTemplate(filename string) *template.Template {
	t, parseErr := template.ParseFiles("template/" + filename)
	if parseErr != nil {
		log.Panicln("NisePostGo: ", parseErr)
	}
	return t
}

func LoadWebContent(filename string) *template.Template {
	t, parseErr := template.ParseFiles("web" + filename)
	if parseErr != nil {
		log.Panicln("NisePostGo: ", parseErr)
	}
	return t
}
