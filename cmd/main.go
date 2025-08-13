package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"sigs.k8s.io/yaml"

	"github.com/sarems/pmfp/internal/config"
)

type ProxyHandler struct {
	config.ConfigElement
}

func (p *ProxyHandler) HandleHttp(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: Handling HTTP request for %s", r.URL.String())

	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		log.Printf("ERROR: Could not create proxy request: %v", err)
		return
	}

	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}

	var client *http.Client

	if p.ConfigElement.ProxyServer == nil {
		client = &http.Client{
			Timeout: 3 * time.Second,
		}
	} else {
		transport := &http.Transport{
			Proxy: http.ProxyURL(p.ConfigElement.ProxyServer),
		}
		client = &http.Client{
			Timeout:   3 * time.Second,
			Transport: transport,
		}
	}

	p.ConfigElement.ApplyRequestMiddleware(proxyReq)

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		log.Printf("ERROR: Could not forward request: %v", err)
		return
	}
	defer resp.Body.Close()

	for h, val := range resp.Header {
		w.Header()[h] = val
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to forward proxy config")
	targetPort := flag.Int("port", 8080, "Port number to bind to")

	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var config config.ConfigElement
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		os.Exit(1)
	}

	proxyHandler := ProxyHandler{ConfigElement: config}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", *targetPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxyHandler.HandleHttp(w, r)
		}),
	}

	log.Printf("🚀 Starting forwarding proxy server on port %d...\n", *targetPort)
	log.Printf("Configure your browser or OS to use HTTP Proxy at 127.0.0.1:%d", *targetPort)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("FATAL: Server failed to start: %v", err)
	}
}
