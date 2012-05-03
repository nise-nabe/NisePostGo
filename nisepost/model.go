package nisepost

import (
	"crypto/sha512"
	"hash"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"log"
	"os"
        "io"
        "fmt"
)

var (
	db *mgo.Database
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

type NisePostGoUser struct {
	Username string
	Password string
}

func NewUser(username, password string) *NisePostGoUser {
  return &NisePostGoUser{username, Encrypto(password)}
}

func Authenticate(username, password string) bool {
	user := NisePostGoUser{}
	return db.C("User").Find(bson.M{"username": username, "password": Encrypto(password)}).One(&user) == nil
}

func Encrypto(s string) string {
	var h hash.Hash = sha512.New()
        io.WriteString(h, s)
        return fmt.Sprintf("%x", h.Sum(nil))
}

func (user *NisePostGoUser) Save() {
  db.C("User").Insert(&user)
}
