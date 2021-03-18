package main

import (
	"MyEvents/boocking-api/src/config"
	"MyEvents/boocking-api/src/server"
	"flag"
	"log"
	"net/http"
)

func main() {
	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()

	conf, _ := config.ExtractConfiguration(*confPath)
	s := server.NewServer(conf)
	httpErrChan := make(chan error)
	defer close(httpErrChan)
	httpsErrChan := make(chan error)

	go func() {
		log.Printf("Listening on addr %s\n", conf.RestfulEndpoint)
		httpErrChan <- http.ListenAndServe(conf.RestfulEndpoint, s)
	}()

	go func() {
		log.Printf("Listening on addr %s\n", conf.SecureRestfulEndpoint)
		httpsErrChan <- http.ListenAndServeTLS(conf.SecureRestfulEndpoint, "https/cert.pem", "https/key.pem", s)
	}()

	select {
	case err := <-httpErrChan:
		log.Fatalf("HTTP error: %s", err)
	case err := <-httpsErrChan:
		log.Fatalf("HTTPS error: %s", err)
	}
}
