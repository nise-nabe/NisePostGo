package main

import (
        "launchpad.net/mgo"
        "strconv"
)

type Goblog struct {
        Content string
}

func main() {
        session, _ := mgo.Dial("localhost")
        defer session.Close()
        session.SetMode(mgo.Monotonic, true)
        c := session.DB("test").C("goblog")
        var i int64
        for i = 0; i < 10000000000000; i++ {
                c.Insert(&Goblog{strconv.FormatInt(i, 10)})
        }
}
