package main

import (
	"log"
	"net/http"
)

func main() {
	filePathRoot := "."
	port := "8080"
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("/healthz", handlerReadiness)

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
