package main

import (
	"database/sql"
	"log"
	"net/http"

	internalhandlers "github.com/gabsfranca/gerador-nf-faturamento/internal/internalHandlers"
	"github.com/gabsfranca/gerador-nf-faturamento/internal/repository"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "postgres"
	dbPass := "senha"
	dbName := "invoice_service_db"

	connStr := "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("falha na con com o pg: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("falha no ping do pg: ", err)
	}

	repo := repository.NewPostgresInvoiceRepository(db)
	handler := internalhandlers.NewInoiceHandler(repo)

	router := mux.NewRouter()
	router.HandleFunc("/invoices", handler.CreateInvoice).Methods("POST")
	router.HandleFunc("/invoices/NFs/{NF}", handler.GetInvoiceByNF).Methods("GET")

	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := "8081"
	log.Printf("serviço de notas fiscais iniciado na porta %v", port)
	log.Fatal(http.ListenAndServe(":"+port, corsOptions(router)))
}
