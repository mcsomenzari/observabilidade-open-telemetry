package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mcsomenzari/temperature-system-by-cep-otel/internal/application/dto"
	"github.com/mcsomenzari/temperature-system-by-cep-otel/internal/application/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func getWeatherServiceB(ctx context.Context, city string) (*dto.WeatherResponse, error) {
	s := fmt.Sprintf(`{"city": "%s"}`, city)
	jsonVar := bytes.NewBuffer([]byte(s))

	req, err := http.NewRequestWithContext(ctx, "POST", "http://service-b:8081/v1/weather", jsonVar)
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseWeather dto.WeatherResponse
	err = json.Unmarshal(res, &responseWeather)
	if err != nil {
		return nil, err
	}

	return &responseWeather, nil
}

func (we *WebServer) CreateServerA() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	//r.Use(middleware.Timeout(60 * time.Second))
	router.Handle("/metrics", promhttp.Handler())
	router.Post("/v1/temperature", we.postZipCode)

	return router
}

func (h *WebServer) postZipCode(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := h.OtelData.OTELTracer.Start(ctx, h.OtelData.RequestNameOTEL)
	defer span.End()

	w.Header().Add("Content-Type", "application/json")

	var request dto.RequestViaCep

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	if len(request.ZipCode) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	start := time.Now()
	responseViaCep, err := services.GetViaCepApiService(request.ZipCode)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	log.Println("response ViaCep:", time.Since(start))

	ctx, spanTime := h.OtelData.OTELTracer.Start(ctx, "SPAN_TIME_VIACEP: "+time.Since(start).String())
	spanTime.End()

	if responseViaCep.Erro {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	responseWeather, error := getWeatherServiceB(ctx, responseViaCep.Localidade)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "an error occurred when calling serviceB - "+error.Error(), http.StatusInternalServerError)
		return
	}

	response := services.FormatTemperatureService(responseViaCep.Localidade, responseWeather.Current.TempC)
	//response := responseViaCep

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
