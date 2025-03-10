package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gabsfranca/gerador-nf-estoque/internal/handlers"
	"github.com/gabsfranca/gerador-nf-estoque/internal/repository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "senha")
	dbName := getEnv("DB_NAME", "gerador_nf_estoque")

	connStr := "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Falha na conexao do pg: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Falha no ping do pg: %v", err)
	}

	repo := repository.NewPostgresProductRepository(db)
	handler := handlers.NewProductHandler(repo)

	router := mux.NewRouter()
	router.HandleFunc("/products", handler.Create).Methods("POST")
	router.HandleFunc("/products", handler.List).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", handler.GetById).Methods("GET")
	router.HandleFunc("/products/serialNumber/{serialNumber}", handler.GetBySerialNumber).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}/stock", handler.UpdateStock).Methods("PUT")

	port := getEnv("PORT", "8080")
	log.Printf("serviço de estoque iniciado com sucesso na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
