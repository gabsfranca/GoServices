package internalhandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gabsfranca/gerador-nf-faturamento/internal/domain"
	"github.com/gabsfranca/gerador-nf-faturamento/internal/repository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type InvoiceHandler struct {
	repo              repository.InvoiceRepository
	productServiceURL string
	httpClient        *http.Client
}

func NewInoiceHandler(repo repository.InvoiceRepository) *InvoiceHandler {
	return &InvoiceHandler{
		repo:              repo,
		productServiceURL: "http://localhost:8080",
		httpClient:        &http.Client{},
	}
}

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("funcao createInvoice chamada")
	var invoice domain.Invoice
	err := json.NewDecoder(r.Body).Decode(&invoice)

	if err != nil {
		fmt.Print("erro ao decodigicar nf: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	invoice.Status = "aberta"

	if invoice.Type == "OUT" {
		for _, item := range invoice.Products {

			url := fmt.Sprintf("%s/products/serialNumber/%s", h.productServiceURL, item.SerialNumber)

			res, err := h.httpClient.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if res.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(res.Body)
				http.Error(w, string(body), res.StatusCode)
				return
			}

			var product struct {
				CurrentStock int `json:"currentStock"`
			}

			if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if product.CurrentStock < int(item.Quantity) {
				http.Error(
					w,
					fmt.Sprintf("estoque insuficiente para o produto %s", item.SerialNumber),
					http.StatusBadRequest,
				)
				return
			}
		}
	}

	if err := h.repo.CreateInvoice(r.Context(), &invoice); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if invoice.Nf == "" || (invoice.Type != "IN" && invoice.Type != "OUT") {
		http.Error(w, "campos invalidos", http.StatusBadRequest)
		return
	}

	if invoice.Status == "" {
		invoice.Status = "Aberta"
	}

	err = h.repo.CreateInvoice(r.Context(), &invoice)
	if err != nil {
		fmt.Println("erro ao salvar no banco: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if invoice.Type == "IN" || invoice.Type == "OUT" {
		success := true
		var adjustments []struct {
			SerialNumber string
			Delta        int
		}

		for _, item := range invoice.Products {
			delta := int(item.Quantity)
			if invoice.Type == "OUT" {
				delta = -delta
			}

			url := fmt.Sprintf("%s/products/%s/update-stock", h.productServiceURL, item.SerialNumber)
			fmt.Println("tentando mandar para: ", url)

			requestBody := struct {
				Delta int `json:"delta"`
			}{Delta: delta}

			jsonBody, _ := json.Marshal(requestBody)
			req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := h.httpClient.Do(req)
			if err != nil {
				success = false
				break
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				success = false
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("falha no update do estoque: %s\n", string(body))
				break
			}

			adjustments = append(adjustments, struct {
				SerialNumber string
				Delta        int
			}{SerialNumber: item.SerialNumber, Delta: delta})
		}

		if !success {
			for _, adj := range adjustments {
				revertDelta := -adj.Delta
				url := fmt.Sprintf("%s/products/%s/update-stock", h.productServiceURL, adj.SerialNumber)
				requestBody := struct {
					Delta int `json:"delta"`
				}{Delta: revertDelta}

				jsonBody, _ := json.Marshal(requestBody)
				req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
				h.httpClient.Do(req)
			}
		} else {
			invoice.Status = "fechada"
			if err := h.repo.UpdateInvoiceStatus(r.Context(), invoice.Id, invoice.Status); err != nil {
				fmt.Printf("erro ao atualizar status: %v\n", err)
			}
		}
	}

	w.Header().Set("Content-Type", "applicfation/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Println("criando invoice: ", invoice)
	json.NewEncoder(w).Encode(invoice)
}

func (h *InvoiceHandler) GetInvoiceByNF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nf := vars["nf"]

	product, err := h.repo.GetInvoiceByNF(r.Context(), nf)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
