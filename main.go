package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func client(clientName string) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(
		"./certs/"+clientName+".pem", "./certs/"+clientName+"-key.pem")
	if err != nil {
		log.Println("CLEINT ERROR:", err)
		return
	}

	// TRUSTING THE SERVER, CLIENT AUTH ONLY
	// Load CA cert
	// caCert, err := ioutil.ReadFile(*caFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// TRUSTING THE SERVER, CLIENT AUTH ONLY
		// RootCAs:      caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	resp, err := client.Get("https://localhost:8443/echo")
	if err != nil {
		log.Println("CLEINT ERROR:", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("CLEINT ERROR:", err)
		return
	}
	log.Println(string(data))
}

func server() error {
	caCert, err := ioutil.ReadFile("./certs/ca.pem")
	if err != nil {
		log.Fatalf("failed to load cert: %s", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(
		"./certs/server.pem", "./certs/server-key.pem")

	tlsConfig := &tls.Config{
		// server certificate which is _not_ validated by the client, but used
		// for encryption
		Certificates: []tls.Certificate{cert},
		// used to verify the client cert is signed by the CA
		ClientCAs: caCertPool,
		// this requires a valid client cert to be supplied during handshake
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      "localhost:8443",
		TLSConfig: tlsConfig,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		res := make(map[string]interface{})
		res["url"] = req.URL.String()
		res["remoteAddr"] = req.RemoteAddr

		if len(req.TLS.PeerCertificates) > 0 {
			res["client"] = req.TLS.PeerCertificates[0].Subject.CommonName
		}
		json.NewEncoder(w).Encode(res)
	})

	http.Handle("/echo", handler)

	// listen using the server certificate which is validated by the client
	return server.ListenAndServeTLS(
		"./certs/server.pem", "./certs/server-key.pem")
}

func main() {
	log.Println("starting server")
	go server()
	time.Sleep(100 * time.Millisecond)
	log.Println("running client")
	client("service-1234@accounts.example.com")
	client("service-3456@accounts.example.com")
	client("service-4567@accounts.example.com")
	time.Sleep(5 * time.Second)
}
