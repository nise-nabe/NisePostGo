package nisepost

import (
	"code.google.com/p/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
)

type NisePostGo struct {
	Content string
}

type NisePostGoHandler struct {
	handler func(http.ResponseWriter, *http.Request)
}

func (h *NisePostGoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	h.handler(w, r)
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
	http.Handle("/", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t := LoadTemplate(w, "web"+r.URL.Path)
			t.Execute(w, nil)
			return
		}
		session, _ := store.Get(r, "session")
		log.Println(session.Values["name"])
		log.Println(session.Values["role"])
		t := LoadTemplate(w, "template/index.html")
		t.Execute(w, nil)
	}})
	http.Handle("/login", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		log.Println(session.Values)
		if session.Values["role"] == nil {
			session, _ = store.New(r, "session")
			session.Values["role"] = "Anonymous"
		} else if session.Values["role"] != "Anonymous" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		session.Save(r, w)
		t := LoadTemplate(w, "template/login.html")
		t.Execute(w, session)
	}})
	http.Handle("/login/check", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		switch r.Method {
		case "POST":
			if session.Values["role"] != "Anonymous" {
				break
			}
			username, password := r.FormValue("username"), r.FormValue("password")
			if !Authenticate(username, password) {
				break
			}
			session.Values["name"] = username
			session.Values["role"] = "User"
			session.Save(r, w)
			log.Println("User Authorized")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		session.Values["_flash"] = make([]interface{}, 0)
		session.AddFlash("login was not succeeded!")
		session.Save(r, w)
		log.Println("User Unauthorized")
		http.Redirect(w, r, "/login", http.StatusFound)
	}})
}

func LoadTemplate(w http.ResponseWriter, filename string) *template.Template {
	t, parseErr := template.ParseFiles(filename)
	if parseErr != nil {
		log.Panicln("NisePostGo: ", parseErr)
	}
	return t
}
