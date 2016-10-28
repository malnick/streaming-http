package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	closed := w.(http.CloseNotifier).CloseNotify()
	last_msg := "unset"
	flusher, _ := w.(http.Flusher)
	for {
		select {
		case <-closed:
			fmt.Println("Client closed connection")
			fmt.Printf("Last Message: %s\n", last_msg)
			return
		default:
			t, _ := time.Now().MarshalJSON()
			last_msg = string(t)
			copy(t, "\n")
			w.Write([]byte(t))

			flusher.Flush()
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/test", TestHandler)
	fmt.Println("Running tester on :8000")
	http.ListenAndServe(":8000", r)
}
