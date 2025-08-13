package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytionbr/temperature_system/service"
)

func TestWeatherHandler_Success(t *testing.T) {
	cepSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"localidade":"São Paulo","erro":"false"}`)
	}))
	defer cepSrv.Close()
	service.CepAPIBaseURL = cepSrv.URL

	weatherSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"current":{"temp_c":20.0}}`)
	}))
	defer weatherSrv.Close()
	WeatherAPIBaseURL = weatherSrv.URL

	t.Setenv("WEATHER_API_KEY", "dummy")

	req := httptest.NewRequest("GET", "/weather?cep=01001000", nil)
	rec := httptest.NewRecorder()
	WeatherHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("wanted 200, got %d", rec.Code)
	}
	var body map[string]float64
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["temp_C"] != 20.0 {
		t.Fatalf("temp_C wrong: %v", body["temp_C"])
	}
	if body["temp_F"] != 20.0*1.8+32 {
		t.Fatalf("temp_F wrong: %v", body["temp_F"])
	}
	if body["temp_K"] != 20.0+273 {
		t.Fatalf("temp_K wrong: %v", body["temp_K"])
	}
}

func TestWeatherHandler_InvalidCEP(t *testing.T) {
	req := httptest.NewRequest("GET", "/weather?cep=123", nil)
	rec := httptest.NewRecorder()
	WeatherHandler(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("wanted 422, got %d", rec.Code)
	}
}

func TestWeatherHandler_ViaCEP_NotFound(t *testing.T) {
	cepSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"erro":"true"}`)
	}))
	defer cepSrv.Close()
	service.CepAPIBaseURL = cepSrv.URL

	t.Setenv("WEATHER_API_KEY", "dummy")
	WeatherAPIBaseURL = "http://teste"

	req := httptest.NewRequest("GET", "/weather?cep=99999999", nil)
	rec := httptest.NewRecorder()
	WeatherHandler(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("wanted 404, got %d", rec.Code)
	}
}

func TestWeatherHandler_NoAPIKey(t *testing.T) {
	cepSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"localidade":"São Paulo","erro":"false"}`)
	}))
	defer cepSrv.Close()
	service.CepAPIBaseURL = cepSrv.URL

	req := httptest.NewRequest("GET", "/weather?cep=01001000", nil)
	rec := httptest.NewRecorder()
	WeatherHandler(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("wanted 500, got %d", rec.Code)
	}
}
