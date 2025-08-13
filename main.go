package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mytionbr/temperature_system/handler"
)

func main() {
	_ = godotenv.Load()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API to consult the weather by CEP"))
	})

	r.Get("/cep/{cep}", handler.CEPHandler)

	r.Get("/weather", handler.WeatherHandler)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
