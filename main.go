package main

import (
	"code.google.com/p/gorilla/sessions"
	"crypto/sha512"
	"hash"
	"html/template"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"log"
	"net/http"
	"os"
)

type NisePostGoHandler struct {
	handler func(http.ResponseWriter, *http.Request)
}

func (h *NisePostGoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	h.handler(w, r)
}

func init() {
	store = sessions.NewCookieStore([]byte("NiseGoPostSecret"))
	initDB()
	initRouting()
}

var (
	store *sessions.CookieStore
	db    *mgo.Database
)

func initDB() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Panicln("NisePostGo: ", err)
		os.Exit(1)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB("NisePostGo")
}

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
			var h hash.Hash = sha512.New()
			h.Write([]byte(password))
			user := NisePostGoUser{}
			err := db.C("User").Find(bson.M{"Username": username, "Password": h.Sum(nil)}).One(&user)
			if err != nil {
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

func main() {
	log.Println("start")
	if err := http.ListenAndServe(":25565", nil); err != nil {
		log.Fatal("NisePostGo: ", err)
	}
}

func LoadTemplate(w http.ResponseWriter, filename string) *template.Template {
	t, parseErr := template.ParseFiles(filename)
	if parseErr != nil {
		log.Panicln("NisePostGo: ", parseErr)
	}
	return t
}

type NisePostGo struct {
	Content string
}

type NisePostGoUser struct {
	Username string
	Password string
}
