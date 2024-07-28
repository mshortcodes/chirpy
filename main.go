package main

import (
	"chirpy/internal/database"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
}

func main() {
	filePathRoot := "."
	port := "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()

	fsHandler := http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot))))
	mux.Handle("/app/*", fsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Fatal(srv.ListenAndServe())
}
