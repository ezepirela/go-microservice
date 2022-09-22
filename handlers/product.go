package handlers

import (
	"log"
	"net/http"

	"github.com/ezepirela/go-microservice/data"
)

type Product struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) GetProducts(w http.ResponseWriter, r *http.Request) {
	productList := data.GetProducts()

	err := productList.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusInternalServerError)
	}
}
