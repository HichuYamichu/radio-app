package app

import (
	"fmt"
	"net/http"
	"time"
)

func NewHandler() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(serve))
	return router
}

func serve(w http.ResponseWriter, r *http.Request) {
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
			w.Write(c.Value())
			flusher.Flush()
			time.Sleep(time.Second)
		}
	}
}
