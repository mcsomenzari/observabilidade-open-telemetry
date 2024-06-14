package services

import "github.com/mcsomenzari/temperature-system-by-cep-otel/internal/application/dto"

func FormatTemperatureService(localidade string, celsius float64) *dto.TemperatureResponse {
	return &dto.TemperatureResponse{
		City:       localidade,
		Celsius:    celsius,
		Kelvin:     celsius + 273,
		Fahrenheit: (celsius * 1.8) + 32,
	}
}
