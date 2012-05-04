package nisepost

import (
	"github.com/nise-nabe/NisePostGo/nisepost/sessions"
	"html/template"
	"log"
	"net/http"
)

type NisePostGoHandler struct {
	handler func(http.ResponseWriter, *http.Request, *sessions.Session)
}

type NisePostGoPageBean struct {
	HasError bool
	Errors   []interface{}
}

func (h *NisePostGoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	h.handler(w, r, sessions.Get(r))
}

func Init() {
	initDB()
	initRouting()
}

var (
	tmpl = initTemplate()
)

func initRouting() {
	http.Handle("/", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		if r.URL.Path != "/" {
			t := loadWebContent(r.URL.Path)
			t.Execute(w, nil)
			return
		}
		log.Println(s.Values["name"])
		log.Println(s.Values["role"])
		tmpl.ExecuteTemplate(w, "index", nil)
	}})
	http.Handle("/login", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		if s.Values["role"] == nil {
			s = sessions.New(r)
			s.Values["role"] = "Anonymous"
		} else if s.Values["role"] != "Anonymous" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		log.Println(s.Values["hasError"])
		s.Save(r, w)
		tmpl.ExecuteTemplate(w, "login", NisePostGoPageBean{s.Values["hasError"].(bool), s.Flashes()})
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
		s = sessions.New(r)
		s.Values["role"] = "Anonymous"
		s.Values["_flash"] = make([]interface{}, 0)
		s.AddFlash("logout was succeeded!")
		s.Values["hasError"] = true
		s.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusFound)
	}})
	http.Handle("/register", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		switch r.Method {
		case "GET":
			if s.Values["role"] != "Anonymous" {
				break
			}
			tmpl.ExecuteTemplate(w, "register", s)
		}
	}})
	http.Handle("/register/post", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
		if s.Values["role"] != "Anonymous" {
			return
		}
		switch r.Method {
		case "POST":
			username := r.FormValue("username")
			if IsExistUser(username) {
				s.AddFlash("the username was exist")
				http.Redirect(w, r, "/register", http.StatusFound)
			}
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

func initTemplate() *template.Template {
	t, parseErr := template.ParseGlob("template/*.tmpl")
	if parseErr != nil {
		log.Panicln("NisePostGo: ", parseErr)
	}
	return t
}

func loadWebContent(filename string) *template.Template {
	t, parseErr := template.ParseFiles("web" + filename)
	if parseErr != nil {
		log.Panicln("NisePostGo: ", parseErr)
	}
	return t
}
