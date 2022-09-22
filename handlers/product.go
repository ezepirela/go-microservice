package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	if r.Method == http.MethodPost {
		p.AddProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("PUT")
		path := r.URL.Path

		regex := regexp.MustCompile(`/([0-9]+)`)
		groupIDs := regex.FindAllStringSubmatch(path, -1)

		if len(groupIDs) != 1 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(groupIDs[0]) != 2 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := groupIDs[0][1]
		id, _ := strconv.Atoi(idString)

		p.UpdateProduct(id, w, r)
		// p.l.Println("id", id)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) UpdateProduct(id int, w http.ResponseWriter, r *http.Request) {
	product := &data.Product{}

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to parse body", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, product)

	if err == data.ErrProductNotFound {
		http.Error(w, data.ErrProductNotFound.Error(), http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (p *Product) GetProducts(w http.ResponseWriter, r *http.Request) {
	productList := data.GetProducts()

	err := productList.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusInternalServerError)
	}
}

func (p *Product) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST route")

	product := &data.Product{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to parse body to JSON", http.StatusInternalServerError)
	}

	data.AddProduct(product)
	p.l.Printf("Prod: %#v", product)
}
