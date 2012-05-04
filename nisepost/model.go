package nisepost

import (
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
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

type User struct {
	Username string
	Password string
}

func NewUser(username, password string) *User {
	return &User{username, Encrypto(password)}
}

func IsExistUser(username string) bool {
	user := User{}
	return db.C("user").Find(bson.M{"username": username}).One(&user) == nil
}

func Authenticate(username, password string) bool {
	user := User{}
	return db.C("User").Find(bson.M{"username": username, "password": Encrypto(password)}).One(&user) == nil
}

func Encrypto(s string) string {
	var h hash.Hash = sha512.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (user *User) Save() {
	db.C("User").Insert(&user)
}
