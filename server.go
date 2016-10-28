package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	closed := w.(http.CloseNotifier).CloseNotify()
	for {
		select {
		case <-closed:
			fmt.Println("Client closed connection")
			return
		default:
			t, _ := time.Now().MarshalJSON()
			fmt.Println("Sending: %s", t)
			copy(t, "\n")
			w.Write([]byte(t))
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/test", TestHandler)
	r.Headers("Transfer-Encoding", "chunked")
	r.Headers("Connection", "keep-alive")
	fmt.Println("Running tester on :8000")
	http.ListenAndServe(":8000", r)
}
