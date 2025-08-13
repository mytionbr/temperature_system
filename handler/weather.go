package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/mytionbr/temperature_system/model"
	"github.com/mytionbr/temperature_system/service"
	"github.com/mytionbr/temperature_system/utils"
)

var WeatherAPIBaseURL = "http://api.weatherapi.com/v1/current.json"

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cep := r.URL.Query().Get("cep")

	cep, err := service.CepValitation(cep)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		newError := utils.NewError(err, http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(newError)
		return
	}

	via, err := service.SearchLocationByCEP(cep)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		newError := utils.NewError(err, http.StatusNotFound)
		json.NewEncoder(w).Encode(newError)
		return
	}

	isErro := via.Erro == "true"

	if isErro {
		w.WriteHeader(http.StatusNotFound)
		newError := utils.NewError(errors.New("can not find zipcode"), http.StatusNotFound)
		json.NewEncoder(w).Encode(newError)
		return
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		newError := utils.NewError(errors.New("server misconfiguration"), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newError)
		return
	}

	city := via.Localidade
	cityEscaped := url.QueryEscape(city)
	reqURL := fmt.Sprintf("%s?key=%s&q=%s", WeatherAPIBaseURL, apiKey, cityEscaped)

	resp, err := http.Get(reqURL)

	if err != nil || resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		newError := utils.NewError(errors.New("error fetching weather data"), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newError)
		return
	}

	defer resp.Body.Close()

	var wr model.WeatherApiResp

	if err := json.NewDecoder(resp.Body).Decode(&wr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		newError := utils.NewError(errors.New("error parsing weather response"), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newError)
		return
	}

	c := wr.Current.TempC
	f := c*1.8 + 32
	k := c + 273

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]float64{
		"temp_C": float64(c),
		"temp_F": float64(f),
		"temp_K": float64(k),
	})
}
