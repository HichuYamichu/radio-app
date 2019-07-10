package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var buffer = make([]byte, 40000)
var chunk []byte

func main() {

	f, err := os.Open("./test.mp3")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	go loadBuffer(f)

	http.HandleFunc("/audio", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}

		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Content-Type", "audio/wav")
		for {
			w.Write(chunk)
			flusher.Flush()
			time.Sleep(time.Second)
		}
	})

	flag.Parse()
	log.SetFlags(0)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func loadBuffer(f *os.File) {
	throttle := time.Tick(time.Second)
	for {
		<-throttle
		t := time.Now()
		fmt.Println(t.Format("20060102150405"))
		bytesread, err := f.Read(buffer)
		chunk = buffer[:bytesread]

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		}
	}
}
