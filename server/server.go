package server

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	pipeReader, pipeWriter = io.Pipe()
	inChan                 = make(chan string, 10)
)

func StdOutHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request made on /stdout")

	flusher, _ := w.(http.Flusher)
	for {
		select {
		case fromIn := <-inChan:
			withTags := fmt.Sprintf("FROM_SERVER[%s]", fromIn)
			log.Infof("ECHO /stdin -> /stdout: %s", withTags)
			w.Write([]byte(withTags))
			flusher.Flush()
		default:
			log.Warn("Nothing received on /stdin, waiting...")
			time.Sleep(1 * time.Second)
		}
	}
}

func StdInHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Request made on /stdin\n%+v", r)

	if r.TransferEncoding[0] != "chunked" {
		log.Errorf("ERROR: Transfer encoding not chunked, Got, %s", r.TransferEncoding[0])
		return
	}

	reader := bufio.NewReader(r.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		log.Infof("Message received from client on /stdin: %s", string(line))
		inChan <- string(line)
	}
}

func StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/stdin", StdInHandler)
	r.HandleFunc("/stdout", StdOutHandler)

	fmt.Println("Running tester on :8000")
	http.ListenAndServe(":8000", r)
}
