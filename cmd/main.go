package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/hichuyamichu/radio-app/app"
)

var addr = flag.String("addr", "localhost:3000", "http service address")
var store = flag.String("store", "./store", "path to mp3 storage")

func main() {
	flag.Parse()
	log.SetFlags(0)

	go app.Start(*store)
	handler := app.NewHandler()
	log.Fatal(http.ListenAndServe(*addr, handler))
}
