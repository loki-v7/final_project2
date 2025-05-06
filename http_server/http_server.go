package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net"
	"strings"
	"sync"
	
)

var (
	messages []string
	mutex sync.RWMutex
)



func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		mutex.RLock()
		defer mutex.RUnlock()
		for _, msg := range messages {
			fmt.Fprintf(w, "%s<br>\n", msg)
		}

	}
}

// /live endpoint returns the messages as JSON
func live_update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.RLock()
		defer mutex.RUnlock()

		// Convert messages to JSON and send as response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

func server_static_files(){
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./public")))) // Serves index.html and other static files
}
func run_server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	http.HandleFunc("/live", live_update())
	fmt.Println("Server starting on port 9001...")

	log.Fatal(http.ListenAndServe("0.0.0.0:9001", nil))
}

func get_words(hostname string) {
	server, err := net.ResolveUDPAddr("udp4", hostname)

	if err != nil {
		fmt.Println("Server: (WARNING/ERROR) could not resolve/reconize hostname -> ", err)
		return
	
	} 

	fmt.Println("Server: (Success) resolved hostname")
	
	connect, err := net.ListenUDP("udp4", server)
	if err != nil {
		fmt.Println("Server: (WARNING/ERROR) could not listen on parsed")
		return
	
	}

	fmt.Println("Server: (Success) Server listening on the parsed address ")	
	defer connect.Close()
	buffer := make([]byte, 1024)

	for {
		network, _, err := connect.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Server: (ERROR) Failed to read from UDP ->", err)
			return // Exit on error
		}

		msg := string(buffer[:network])
		fmt.Println("Server: (Success) Got string payload ->", msg)

		// Locking for thread safety when appending the message
		mutex.Lock()
		messages = append(messages, msg)
		mutex.Unlock()

		// If we receive "ServerStop", shut down the server
		if strings.TrimSpace(msg) == "ServerStop" {
			fmt.Println("Server: Shutting down server...")
			return
		}
	}
}

func main() {
	go get_words("0.0.0.0:9002")
	run_server()
}
