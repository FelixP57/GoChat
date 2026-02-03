package main

import (
	"log"
	"net/http"
	"flag"
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	// parse addr
	flag.Parse()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
	})

	// create root ctx and cancelfunc to cancel retention map goroutine
	rootCtx := context.Background()
	ctx, cancel := context.WithCancel(rootCtx)

	defer cancel()

	err = setupAPI(ctx, mux)
	if err != nil {
		log.Fatal("Error setting up the API: ", err)
		return
	}

	// serve on designated addr
	// err = http.ListenAndServeTLS(*addr, "localhost+2.pem", "localhost+2-key.pem", c.Handler(mux))
	err = http.ListenAndServe(*addr, c.Handler(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}

// start all routes and associated handlers
func setupAPI(ctx context.Context, mux *http.ServeMux) error {
	// hub to handle websocket connections
	hub, err := newHub(ctx)
	if err != nil {
		return err
	}
	go hub.run()


	// serve the ./frontend dir at route /
	/*
	fs := http.FileServer(http.Dir("../frontend/"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHome(hub, w, r)
	})
	*/
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.serveWs(w, r)
	})
	mux.HandleFunc("/login", hub.loginHandler)
	mux.HandleFunc("/signup", hub.signupHandler)
	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, len(hub.clients))
	})
	// http.Handle("/frontend/", http.StripPrefix("/frontend", fs))

	return nil
}

