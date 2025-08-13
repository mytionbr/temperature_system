package service

import (
	"errors"
	"regexp"
	"strings"
)

func CepValitation(cep string) (string, error) {
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
