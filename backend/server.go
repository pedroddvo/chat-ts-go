package main

import (
	"log"
	"net/http"
)

const port = "8080"

func main() {
	db := DbNew()
	go db.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(db, w, r)
	})

	log.Println("Initializing server on port " + port)
	err := http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		log.Fatal("Failed to initialize server: ", err)
	}
}
