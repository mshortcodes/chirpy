package main

import (
	"chirpy/internal/database"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	godotenv.Load(".env")

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		log.Println("Debug mode enabled. Deleting database...")
		err := os.Remove("database.json")
		if err != nil {
			log.Println("Couldn't delete database.")
		}
	}

	filePathRoot := "."
	port := "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
	}

	mux := http.NewServeMux()

	fsHandler := http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot))))
	mux.Handle("/app/*", fsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Fatal(srv.ListenAndServe())
}
