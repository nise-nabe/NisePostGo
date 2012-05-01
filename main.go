package main

import (
	"github.com/nise-nabe/NisePostGo/nisepost"
	"log"
	"net/http"
)

func main() {
	log.Println("start")
	nisepost.Init()
	if err := http.ListenAndServe(":25565", nil); err != nil {
		log.Fatal("NisePostGo: ", err)
	}
}
