package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strconv"
)

func InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/products", func(r chi.Router) {
		r.Get("/", GetProducts)
		r.Get("/{id}", GetProduct)
	})

	router.Route("/categories", func(r chi.Router) {
		r.Get("/", GetCategories)
		r.Get("/{id}", GetCategory)
		r.Get("/{id}/*", GetCategory)
	})

	return router
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var (
		products []bson.M
		err      error
	)

	urlParams := r.URL.Query()
	params := make(map[string]string)

	for key, value := range urlParams {
		if len(value) > 0 {
			params[key] = value[0]
		}
	}

	log.Println(params)

	if len(params) > 0 {
		products, err = GetProductsByQuery(params)
	} else {
		products, err = GetAllProducts()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	product, err := GetProductById(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(product)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	type Categories struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	categories := []Categories{
		{
			ID:   1,
			Name: "филимоновская игрушка",
		},
		{
			ID:   2,
			Name: "пастила",
		},
		{
			ID:   3,
			Name: "тульские самовары",
		},
		{
			ID:   4,
			Name: "тульские пряники",
		},
	}

	jsonResponse, err := json.Marshal(categories)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	limitParam := r.URL.Query().Get("limit")

	var limit int
	if limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	category, err := GetProductsByCategory(id, limit)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(category)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
