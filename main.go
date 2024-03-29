package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	WebPort          = 8080
	InternalWebPort  = 8081
	interruptSignals = []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	}
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	waitGroup, ctx := errgroup.WithContext(ctx)
	runInternalServer(waitGroup, ctx)
	runMainServer(waitGroup, ctx)
	err := waitGroup.Wait()
	if err != nil {
		fmt.Printf("Error during wait: %v\n", err)
		os.Exit(1)
	}
}

func runInternalServer(waitGroup *errgroup.Group, ctx context.Context) {
	internalServer := newInternalServer()
	waitGroup.Go(func() error {
		fmt.Printf("Starting internal http server on port : %d\n", InternalWebPort)
		if err := internalServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			fmt.Printf("Failed to start internal http server: %v\n", err)
			return err
		}
		return nil

	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		fmt.Println("Gracefully shutting down internal http server")
		err := internalServer.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("Failed to shutdown internal http server: %v\n", err)
			return err
		}
		fmt.Println("Internal http server shutdown complete")
		return nil
	})
}

func runMainServer(waitGroup *errgroup.Group, ctx context.Context) {
	server := newServer()
	waitGroup.Go(func() error {
		fmt.Printf("Starting main http server  on port : %d\n", WebPort)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			fmt.Printf("Failed to start main http server: %v\n", err)
			return err
		}
		return nil
	})
	waitGroup.Go(func() error {
		<-ctx.Done()
		fmt.Println("Gracefully shutting down main http server")
		err := server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("Failed to shutdown main http server: %v\n", err)
			return err
		}
		fmt.Println("Main http server shutdown complete")
		return nil
	})

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
