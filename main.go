// main.go
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

var quotes []Quote

func loadQuotes(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			quote := Quote{
				Text:   strings.TrimSpace(parts[0]),
				Author: strings.TrimSpace(parts[1]),
			}
			quotes = append(quotes, quote)
		} else {
			log.Printf("Skipping malformed line: %s", line)
		}
	}

	return scanner.Err()
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func getQuoteHandler(w http.ResponseWriter, r *http.Request) {

	if len(quotes) == 0 {
		http.Error(w, "No quotes available", http.StatusInternalServerError)
		return
	}

	randomQuote := quotes[random.Intn(len(quotes))]

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(randomQuote)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getLandingPage(w http.ResponseWriter, r *http.Request) {
	welcomeMessage := "Welcome to the Quote of the day application!"

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(welcomeMessage)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "ok"}`)
}

func main() {
	if err := loadQuotes("quotes.txt"); err != nil {
		log.Fatalf("Failed to load quotes: %v", err)
	}
	log.Printf("✅ Successfully loaded %d quotes from file.", len(quotes))

	mux := http.NewServeMux()
	mux.HandleFunc("/", getLandingPage)
	mux.HandleFunc("/quote", getQuoteHandler)
	mux.HandleFunc("/healthz", healthzHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("🚀 Server starting on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	<-stopChan
	log.Println("🔌 Server shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}

	log.Println("✅ Server exited properly")
}
