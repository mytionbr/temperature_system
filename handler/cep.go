package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mytionbr/temperature_system/model"
	"github.com/mytionbr/temperature_system/service"
	"github.com/mytionbr/temperature_system/utils"
)

func CEPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cep := chi.URLParam(r, "cep")

	cep, err := service.CepValitation(cep)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		newError := utils.NewError(err, http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(newError)
		return
	}

	var response model.CepApiResponseData

	response, err = service.SearchLocationByCEP(cep)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		newError := utils.NewError(err, http.StatusNotFound)
		json.NewEncoder(w).Encode(newError)
		return
	}

	isErro := response.Erro == "true"

	if isErro {
		w.WriteHeader(http.StatusUnprocessableEntity)
		newError := utils.NewError(errors.New("invalid zipcode. The zipcode must be 8 digits long"), http.StatusBadRequest)
		json.NewEncoder(w).Encode(newError)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

type CepApiResponse map[string]model.CepApiResponseData
