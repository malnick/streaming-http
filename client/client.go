package client

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

type IOData struct {
	Foo string `json:"foo"`
}

func getStdinRequester(pipeReader *io.PipeReader, host, port, stdinEndpoint string) (*http.Request, error) {
	thisHost := fmt.Sprintf("%s:%s", host, port)
	hitme := url.URL{
		Scheme: "http",
		Host:   thisHost,
		Path:   stdinEndpoint,
	}

	headers := map[string]string{
		"Transfer-Encoding": "chunked",
	}

	req, err := http.NewRequest("POST", hitme.String(), ioutil.NopCloser(pipeReader))
	if err != nil {
		return req, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return req, err
}

func getStdoutRequester(host, port, stdoutEndpoint string) (*http.Request, error) {
	thisHost := fmt.Sprintf("%s:%s", host, port)
	hitme := url.URL{
		Scheme: "http",
		Host:   thisHost,
		Path:   stdoutEndpoint,
	}
	return http.NewRequest("GET", hitme.String(), nil)
}

func startClientGenerator(pipeWriter *io.PipeWriter) {
	for x := 0; x <= 5; x++ {
		time.Sleep(1 * time.Second)
		sendme := fmt.Sprintf("It is now %v\n", time.Now())
		log.Info("SENDING to /stdin %s", sendme)
		fmt.Fprintf(pipeWriter, sendme)
	}
	// Send EOF
	log.Warn("Sending EOF to /stdin")
	fmt.Fprintf(pipeWriter, "")

}

func sendToServer(req *http.Request) {
	log.Infof("Opening stream to %s", req.URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	n, err := io.Copy(os.Stdout, resp.Body)
	log.Fatalf("copied %d, %v", n, err)
}

func readFromServer(req *http.Request) {
	log.Infof("Opening stream to %s", req.URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		log.Infof("ECHO RECEIVED: %s", string(line))
	}
}

func StartClient() {
	cf := flag.NewFlagSet("", flag.ContinueOnError)
	var (
		host           = cf.String("host", "localhost", "IP address of host")
		stdinEndpoint  = cf.String("stdin", "/stdin", "Stdin endpoint on remote")
		stdoutEndpoint = cf.String("stdout", "/stdout", "Stdout endpoint on remote")
		port           = cf.String("port", "8000", "Port to hit on remote")
	)
	cf.Parse(os.Args[2:])

	pipeReader, pipeWriter := io.Pipe()
	stdinReq, err := getStdinRequester(pipeReader, *host, *port, *stdinEndpoint)
	if err != nil {
		panic(err)
	}

	stdoutReq, err := getStdoutRequester(*host, *port, *stdoutEndpoint)
	if err != nil {
		panic(err)
	}

	go startClientGenerator(pipeWriter)
	go sendToServer(stdinReq)
	readFromServer(stdoutReq)

}
