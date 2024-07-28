package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var cep_api = "https://viacep.com.br/ws/%v/json"

func main() {
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

		if response.Erro == "true" {
			w.WriteHeader(http.StatusBadRequest)
			newError := newError(errors.New("cep inválido. O cep precisa ter 8 dígitos"), http.StatusBadRequest)
			json.NewEncoder(w).Encode(newError)
			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)

	})

	http.ListenAndServe(":3000", r)
}

func cepValitation(cep string) (string, error) {
	cep = strings.Replace(cep, " ", "", -1)
	cep = strings.Replace(cep, "-", "", -1)

	if len(cep) != 8 {
		return "", errors.New("cep inválido. O cep precisa ter 8 dígitos")
	}

	isNumeric := regexp.MustCompile(`^[0-9]+$`).MatchString(cep)

	if !isNumeric {
		return "", errors.New("cep inválido. O cep deve ter apenas dígitos numéricos")
	}

	return cep, nil
}

func searchLocationByCEP(cep string) (CepApiResponseData, error) {
	res, err := http.Get(fmt.Sprintf(cep_api, cep))

	if err != nil {
		log.Fatalln(err)
		return CepApiResponseData{}, errors.New("não foi possível verificar a localidade do cep")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
		return CepApiResponseData{}, errors.New("não foi possível verificar a localidade do cep")
	}

	var apiResponse CepApiResponseData
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		log.Fatalln(err)
		return CepApiResponseData{}, errors.New("não foi possível verificar a localidade do cep")
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
	Erro        string `json:"erro,omitempty"`
}
