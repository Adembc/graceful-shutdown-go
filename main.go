package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	WebPort         = 8080
	InternalWebPort = 8081
)

func main() {

	go runInternalServer()
	runMainServer()
}

func runInternalServer() {
	internalServer := newInternalServer()

	fmt.Printf("Starting internal http server on port : %d\n", InternalWebPort)
	if err := internalServer.ListenAndServe(); err != nil {
		fmt.Printf("Failed to start internal http server: %v\n", err)
		os.Exit(1)
	}
}

func runMainServer() {
	server := newServer()

	fmt.Printf("Starting main http server  on port : %d\n", WebPort)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to start main http server: %v\n", err)
		os.Exit(1)
	}

}

func newServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/readyz", readinessHandler)
	mux.HandleFunc("/calculation", calculationHandler)
	return &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", WebPort),
	}
}

func newInternalServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/readyz", readinessHandler)
	mux.HandleFunc("/encryption", encryptionHandler)
	return &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", InternalWebPort),
	}
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, World!\n")
}

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Server is ready!\n")
}

func calculationHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Calculating...")
	// Simulate a long calculation
	time.Sleep(10 * time.Second)
	fmt.Println("Calculation is done!")
	fmt.Fprintf(w, "Calculation is done!\n")
}

func encryptionHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Encrypting ...")
	// Simulate a long calculation
	time.Sleep(10 * time.Second)
	fmt.Println("encrypting is done!")
	fmt.Fprintf(w, "SECRET_KEY!\n")
}
