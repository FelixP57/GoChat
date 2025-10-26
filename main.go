package main

import (
	"log"
	"net/http"
	"flag"
	"context"
	"fmt"

	"github.com/joho/godotenv"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return 
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return 
	}

	http.ServeFile(w, r, "./frontend/index.html")
}

func main() {
	// parse addr
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// create root ctx and cancelfunc to cancel retention map goroutine
	rootCtx := context.Background()
	ctx, cancel := context.WithCancel(rootCtx)

	defer cancel()

	setupAPI(ctx)

	// serve on designated addr
	err = http.ListenAndServeTLS(*addr, "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// start all routes and associated handlers
func setupAPI(ctx context.Context) {
	// hub to handle websocket connections
	hub := newHub(ctx)
	go hub.run()

	// serve the ./frontend dir at route /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHome(hub, w, r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.serveWs(w, r)
	})
	http.HandleFunc("/login", hub.loginHandler)
	http.HandleFunc("/signup", hub.signupHandler)
	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, len(hub.clients))
	})
}

