package nisepost

import (
	"crypto/sha512"
	"hash"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"log"
	"os"
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

func Authenticate(username, password string) bool {
	var h hash.Hash = sha512.New()
	h.Write([]byte(password))
	user := NisePostGoUser{}
	return db.C("User").Find(bson.M{"Username": username, "Password": h.Sum(nil)}).One(&user) == nil
}

func (user NisePostGoUser) Save() {
  db.C("User").Insert(&user)
}
