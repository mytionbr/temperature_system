package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mytionbr/temperature_system/model"
)

var cep_api = "https://viacep.com.br/ws/%v/json"

func SearchLocationByCEP(cep string) (model.CepApiResponseData, error) {
	res, err := http.Get(fmt.Sprintf(cep_api, cep))

	if err != nil {
		log.Fatalln(err)
		return model.CepApiResponseData{}, errors.New("can not find zipcode")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
		return model.CepApiResponseData{}, errors.New("can not find zipcode")
	}

	var apiResponse model.CepApiResponseData
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		log.Fatalln(err)
		return model.CepApiResponseData{}, errors.New("can not find zipcode")
	}

	return apiResponse, nil
}
