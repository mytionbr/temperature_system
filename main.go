package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

var cep_api = "https://viacep.com.br/ws/%v/json"

func main() {
	_ = godotenv.Load()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Api para consultar o tempo de acordo com o cep"))
	})
	r.Get("/cep/{cep}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cep := chi.URLParam(r, "cep")

		cep, err := cepValitation(cep)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			newError := newError(err, http.StatusBadRequest)
			json.NewEncoder(w).Encode(newError)
			return
		}

		var response CepApiResponseData

		response, err = searchLocationByCEP(cep)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			newError := newError(err, http.StatusBadRequest)
			json.NewEncoder(w).Encode(newError)
			return
		}

		if response.Erro {
			w.WriteHeader(http.StatusBadRequest)
			newError := newError(errors.New("invalid zipcode. The zipcode must be 8 digits long"), http.StatusBadRequest)
			json.NewEncoder(w).Encode(newError)
			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)

	})

	r.Get("/weather", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cep := r.URL.Query().Get("cep")

		cep, err := cepValitation(cep)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			newError := newError(err, http.StatusBadRequest)
			json.NewEncoder(w).Encode(newError)
			return
		}

		via, err := searchLocationByCEP(cep)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			newError := newError(err, http.StatusBadRequest)
			json.NewEncoder(w).Encode(newError)
			return
		}

		if via.Erro {
			w.WriteHeader(http.StatusNotFound)
			newError := newError(errors.New("zipcode not found"), http.StatusNotFound)
			json.NewEncoder(w).Encode(newError)
		}

		apiKey := os.Getenv("WEATHER_API_KEY")
		if apiKey == "" {
			w.WriteHeader(http.StatusInternalServerError)
			newError := newError(errors.New("server misconfiguration"), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newError)
			return
		}

		city := via.Localidade
		cityEscaped := url.QueryEscape(city)
		reqURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, cityEscaped)

		resp, err := http.Get(reqURL)

		if err != nil || resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			newError := newError(errors.New("error fetching weather data"), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(newError)
			return
		}

		defer resp.Body.Close()

		var wr weatherApiResp

		if err := json.NewDecoder(resp.Body).Decode(&wr); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			newError := newError(errors.New("error parsing weather response"), http.StatusInternalServerError)
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

	})
	fmt.Println("Servidor rodando na porta 3000")
	http.ListenAndServe(":3000", r)
}

func cepValitation(cep string) (string, error) {
	cep = strings.Replace(cep, " ", "", -1)
	cep = strings.Replace(cep, "-", "", -1)

	if len(cep) != 8 {
		return "", errors.New("invalid zipcode. The zipcode must be 8 digits long")
	}

	isNumeric := regexp.MustCompile(`^[0-9]+$`).MatchString(cep)

	if !isNumeric {
		return "", errors.New("invalid zipcode. The zipcode must contain only numeric digits")
	}

	return cep, nil
}

func searchLocationByCEP(cep string) (CepApiResponseData, error) {
	res, err := http.Get(fmt.Sprintf(cep_api, cep))

	if err != nil {
		log.Fatalln(err)
		return CepApiResponseData{}, errors.New("Could not verify zipcode location.")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
		return CepApiResponseData{}, errors.New("Could not verify zipcode location.")
	}

	var apiResponse CepApiResponseData
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		log.Fatalln(err)
		return CepApiResponseData{}, errors.New("Could not verify zipcode location.")
	}

	return apiResponse, nil
}

type StatusError struct {
	Code    int
	Err     error
	Message string
	Caller  string
}

func newError(err error, code int) StatusError {
	pc, _, line, _ := runtime.Caller(2)
	details := runtime.FuncForPC(pc)

	return StatusError{
		Code:    code,
		Err:     err,
		Message: fmt.Sprintf("%v", err),
		Caller:  fmt.Sprintf("%s#%d", details.Name(), line),
	}
}

type Response struct {
	Message string
}

type TemperatureResponse struct {
	Temp_C float32
	Temp_F float32
	Temp_K float32
}

type CepApiResponse map[string]CepApiResponseData

type CepApiResponseData struct {
	Cep         string `json:"cep"`
	Logradoro   string `json:"logradoro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        bool   `json:"erro,omitempty"`
}

type weatherApiResp struct {
	Current struct {
		TempC float32 `json:"temp_c"`
	} `json:"current"`
}
