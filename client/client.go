package client

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type IOData struct {
	Foo string `json:"foo"`
}

var (
	host     = flag.String("host", "localhost", "IP address of host")
	endpoint = flag.String("endpoint", "/test", "Endpoint to hit on remote")
	port     = flag.Int("port", 80, "Port to hit on remote")
)

func getRequester(pipeReader *io.PipeReader) (*http.Request, error) {
	thisHost := fmt.Sprintf("%s:%s", *host, *port)
	hitme := url.URL{
		Scheme: "http",
		Host:   thisHost,
		Path:   *endpoint,
	}

	headers := map[string]string{
		"Content-Type":      "application/json",
		"Transfer-Encoding": "chunked",
		"Connection":        "keep-alive",
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

func startClientGenerator(pipeWriter *io.PipeWriter) {
	for {
		time.Sleep(1 * time.Second)
		fmt.Fprintf(pipeWriter, "It is now %v\n", time.Now())
	}
}

func sendToServer(req *http.Request) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	log.Printf("Got: %#v", resp)
	n, err := io.Copy(os.Stdout, resp.Body)
	log.Fatalf("copied %d, %v", n, err)
}

func main() {
	pipeReader, pipeWriter := io.Pipe()
	req, err := getRequester(pipeReader)
	if err != nil {
		panic(err)
	}

	go startClientGenerator(pipeWriter)
	go sendToServer(req)

}
