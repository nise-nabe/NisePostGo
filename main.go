package main

import (
        "html/template"
        "launchpad.net/mgo"
        //        "launchpad.net/mgo/bson"
        "log"
        "net/http"
)

type GoblogHandler struct {
        handler func(http.ResponseWriter, *http.Request)
}

func (h *GoblogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        log.Println(r.URL.Path)
        h.handler(w, r)
}

func init() {
        http.Handle("/", &GoblogHandler{func(w http.ResponseWriter, r *http.Request) {
                if r.URL.Path != "/" {
                        t := LoadTemplate(w, "web"+r.URL.Path)
                        t.Execute(w, nil)
                        return
                }
                t := LoadTemplate(w, "template/index.html")
                t.Execute(w, nil)
        }})
        http.Handle("/edit", &GoblogHandler{func(w http.ResponseWriter, r *http.Request) {
                log.Println(r.Method)
                switch r.Method {
                case "GET":
                        t := LoadTemplate(w, "template/edit.html")
                        t.Execute(w, nil)
                case "POST":
                        session, err := mgo.Dial("localhost")
                        if err != nil {
                                log.Panicln("Goblog: ", err)
                        }
                        defer session.Close()
                        content := r.FormValue("content")
                        session.SetMode(mgo.Monotonic, true)
                        c := session.DB("test").C("goblog")
                        err = c.Insert(&Goblog{content})
                        if err != nil {
                                log.Panicln("Goblog: ", err)
                        }
                        handler := http.RedirectHandler("/edit", 200)
                        r.Method = "GET"
                        handler.ServeHTTP(w, r)
                }
        }})
        http.Handle("/mongo", &GoblogHandler{func(w http.ResponseWriter, r *http.Request) {
                session, err := mgo.Dial("localhost")
                if err != nil {
                        log.Panicln("Goblog: ", err)
                }
                defer session.Close()
                session.SetMode(mgo.Monotonic, true)
                c := session.DB("test").C("goblog")
                result := []Goblog{}
                err = c.Find(nil).Limit(1000).All(&result)
                if err != nil {
                        log.Println("Goblog: ", err)
                }
                t := LoadTemplate(w, "template/mongo.html")
                t.Execute(w, result)
        }})
}

type Goblog struct {
        Content string
}

func main() {
        log.Println("start")
        if err := http.ListenAndServe(":25565", nil); err != nil {
                log.Fatal("Goblog: ", err)
        }
}

func LoadTemplate(w http.ResponseWriter, filename string) *template.Template {
        t, parseErr := template.ParseFiles(filename)
        if parseErr != nil {
                log.Panicln("Goblog: ", parseErr)
        }
        return t
}
