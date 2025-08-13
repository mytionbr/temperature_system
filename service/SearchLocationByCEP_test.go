package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchLocationByCEP_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"localidade":"São Paulo","erro":"false"}`)
	}))
	defer srv.Close()

	CepAPIBaseURL = srv.URL

	got, err := SearchLocationByCEP("01001000")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if got.Localidade != "São Paulo" {
		t.Fatalf("wanted São Paulo, got %s", got.Localidade)
	}
}

func TestSearchLocationByCEP_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"erro":"true"}`, http.StatusBadRequest)
	}))
	defer srv.Close()

	CepAPIBaseURL = srv.URL

	_, err := SearchLocationByCEP("99999999")
	if err == nil {
		t.Fatal("expected error when ViaCEP does not return 200")
	}
}
