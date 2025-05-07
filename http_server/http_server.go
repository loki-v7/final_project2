package main //makes this the main package

// imports are libraies
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net"
	"strings"
	"sync"
	
)

// makes global variables
var (
	messages []string
	mutex sync.RWMutex
)


// func is how you def a function
func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		mutex.RLock() //mutex is used to protect variables from breaking 
		defer mutex.RUnlock() // defer means wait untill return and then run
		for _, msg := range messages { //for ever msg in recived print to ui
			fmt.Fprintf(w, "%s<br>\n", msg)
		}

	}
}

// /live endpoint returns the messages as JSON
func live_update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { //funcs for setting up http
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
func run_server() { //makes server endpoints and supply html code
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	http.HandleFunc("/live", live_update())
	fmt.Println("Server starting on port 9001...") //print server starting

	log.Fatal(http.ListenAndServe("0.0.0.0:9001", nil)) // server in network on port 9001 
}

func get_words(hostname string) { //this it to get information on port 9002 from udp packets and give to server
	server, err := net.ResolveUDPAddr("udp4", hostname) //find give host

	if err != nil {
		fmt.Println("Server: (WARNING/ERROR) could not resolve/reconize hostname -> ", err)
		return
	
	} 

	fmt.Println("Server: (Success) resolved hostname")
	
	connect, err := net.ListenUDP("udp4", server) //listen with udp
	if err != nil {
		fmt.Println("Server: (WARNING/ERROR) could not listen on parsed")
		return
	
	}

	fmt.Println("Server: (Success) Server listening on the parsed address ")	
	defer connect.Close() //defer to close
	buffer := make([]byte, 1024) // buffer = list of info 

	for { //for loop for the code to get inputs on everything
		network, _, err := connect.ReadFromUDP(buffer) 
		if err != nil {
			fmt.Println("Server: (ERROR) Failed to read from UDP ->", err)
			return // Exit on error
		}

		msg := string(buffer[:network]) // print the payload
		fmt.Println("Server: (Success) Got string payload ->", msg)

		// Locking for thread safety when appending the message
		mutex.Lock()
		messages = append(messages, msg)
		mutex.Unlock()

		// If we receive "ServerStop", shut down the server
		if strings.TrimSpace(msg) == "ServerStop" { // need to get msg from netcat to close so users cant break (unless they know what there doing)
			fmt.Println("Server: Shutting down server...")
			return
		}
	}
}

func main() {
	go get_words("0.0.0.0:9002") // go function() just runs the func in parallel with other funcs
	run_server() // so run the server and then along the get words on port 9002
}
