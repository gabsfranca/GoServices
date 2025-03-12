package internalHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gabsfranca/gerador-nf-estoque/internal/domain"
	"github.com/gabsfranca/gerador-nf-estoque/internal/repository"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {

	var product domain.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	fmt.Println("request: ", r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("erro decodificando o json: ", err)
		return
	}

	err = h.repo.Create(r.Context(), &product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}

func (h *ProductHandler) GetBySerialNumber(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	serialNumber := vars["serialNumber"]

	product, err := h.repo.GetBySerialNumber(r.Context(), serialNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProduts(w http.ResponseWriter, r *http.Request) {

	products, err := h.repo.GetProduts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serialNumber := vars["serialNumber"]

	var request struct {
		Delta int `json:"delta"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.repo.GetBySerialNumber(r.Context(), serialNumber)

	if err != nil {
		http.Error(w, "produto nao encontrado", http.StatusNotFound)
		return
	}

	newStock := product.CurrentStock + request.Delta

	if newStock < 0 {
		http.Error(w, "estoque insuficiente", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateStock(r.Context(), product.ID, newStock); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		NewStock int `json:"newStock"`
	}{NewStock: newStock})
}
