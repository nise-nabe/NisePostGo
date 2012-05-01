package main

import (
	"html/template"
	"launchpad.net/mgo"
	//        "launchpad.net/mgo/bson"
	"code.google.com/p/gorilla/sessions"
	"log"
	"net/http"
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
		t := LoadTemplate(w, "template/index.html")
		t.Execute(w, nil)
	}})
	http.Handle("/login", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			t := LoadTemplate(w, "template/login.html")
			t.Execute(w, nil)
		case "POST":
			username, password := r.FormValue("username"), r.FormValue("password")
			session, _ := store.New(r, "session")
			session.Values["name"] = username
			session.Save(r, w)
			handler := http.RedirectHandler("/", 200)
			r.Method = "GET"
			handler.ServeHTTP(w, r)
		}
	}})
	http.Handle("/edit", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)
		switch r.Method {
		case "GET":
			t := LoadTemplate(w, "template/edit.html")
			t.Execute(w, nil)
		case "POST":
			session, err := mgo.Dial("localhost")
			if err != nil {
				log.Panicln("NisePostGo: ", err)
			}
			defer session.Close()
			content := r.FormValue("content")
			session.SetMode(mgo.Monotonic, true)
			c := session.DB("test").C("goblog")
			err = c.Insert(&NisePostGo{content})
			if err != nil {
				log.Panicln("NisePostGo: ", err)
			}
			handler := http.RedirectHandler("/", 200)
			r.Method = "GET"
			handler.ServeHTTP(w, r)
		}
	}})
	http.Handle("/mongo", &NisePostGoHandler{func(w http.ResponseWriter, r *http.Request) {
		session, err := mgo.Dial("localhost")
		if err != nil {
			log.Panicln("NisePostGo: ", err)
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("test").C("goblog")
		result := []NisePostGo{}
		err = c.Find(nil).Limit(1000).All(&result)
		if err != nil {
			log.Println("NisePostGo: ", err)
		}
		t := LoadTemplate(w, "template/mongo.html")
		t.Execute(w, result)
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
