package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

var addr = flag.String("addr", "localhost:3000", "http service address")
var store = flag.String("store", "./store", "path to mp3 storage")
var buffer = make([]byte, 40000)
var done = make(chan struct{})

type chunk struct {
	mu  sync.RWMutex
	val []byte
}

func (c *chunk) load(f *os.File) {
	throttle := time.Tick(time.Second)
	for {
		<-throttle
		bytesread, err := f.Read(buffer)
		c.mu.Lock()
		c.val = buffer[:bytesread]
		c.mu.Unlock()
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			done <- struct{}{}
			break
		}
	}
}

func (c *chunk) value() []byte {
	c.mu.RLock()
	val := c.val
	c.mu.RUnlock()
	return val
}

var c chunk

func main() {
	flag.Parse()
	log.SetFlags(0)

	go func() {
		files, err := ioutil.ReadDir(*store)
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
			for _, file := range files {
				f, err := os.Open(path.Join(*store, file.Name()))
				defer f.Close()
				if err != nil {
					fmt.Println(err)
					return
				}

				go c.load(f)
				<-done
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Content-Type", "audio/wav")

		for {
			select {
			case <-r.Context().Done():
				break
			default:
				t := time.Now()
				fmt.Println(t.Format("20060102150405"))
				w.Write(c.value())
				flusher.Flush()
				time.Sleep(time.Second)
			}
		}
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
