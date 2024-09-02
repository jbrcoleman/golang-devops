package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	caBytes, err := os.ReadFile("ca.cert")
	if err != nil {
		log.Fatal(err)
	}
	ca := x509.NewCertPool()
	if !ca.AppendCertsFromPEM(caBytes) {
		log.Fatal("ca.cert not valid")
	}

	cert, err := tls.LoadX509KeyPair("client.cert", "client.key")
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      ca,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get("https://go-demo.localtest.me/common-name")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Body (status %d): %s\n", resp.StatusCode, body)
}
