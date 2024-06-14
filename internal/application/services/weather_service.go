package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"github.com/mcsomenzari/temperature-system-by-cep-otel/internal/application/dto"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func GetWeatherApiService(city string) (*dto.WeatherResponse, error) {
	key := "8c5d11c25d764ebdba311342240502"
	c := removeAccents(city)
	c = strings.ReplaceAll(c, " ", "&nbsp;")
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", key, c)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var weatherResponse dto.WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return nil, err
	}

	return &weatherResponse, nil
}
