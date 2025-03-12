package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gabsfranca/gerador-nf-estoque/internal/internalHandlers"
	"github.com/gabsfranca/gerador-nf-estoque/internal/repository"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "senha")
	dbName := getEnv("DB_NAME", "products_service_db")

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
	handler := internalHandlers.NewProductHandler(repo)

	router := mux.NewRouter()
	router.HandleFunc("/products", handler.Create).Methods("POST")
	router.HandleFunc("/products", handler.GetProduts).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", handler.GetProductById).Methods("GET")
	router.HandleFunc("/products/serialNumber/{serialNumber}", handler.GetBySerialNumber).Methods("GET")
	router.HandleFunc("/products/{serialNumber}/update-stock", handler.UpdateStock).Methods("PUT")

	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := getEnv("PORT", "8080")
	log.Printf("serviço de estoque iniciado com sucesso na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsOptions(router)))
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
