package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mcsomenzari/temperature-system-by-cep-otel/internal/application/dto"
)

func GetViaCepApiService(zipCode string) (*dto.ViaCepResponse, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", zipCode)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var viaCepResponse dto.ViaCepResponse
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		return nil, err
	}

	return &viaCepResponse, nil
}
