package web

import (
	"encoding/json"
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

func (we *WebServer) CreateServerB() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	//r.Use(middleware.Timeout(60 * time.Second))
	router.Handle("/metrics", promhttp.Handler())
	router.Post("/v1/weather", we.postWeatherApi)

	return router
}

func (h *WebServer) postWeatherApi(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := h.OtelData.OTELTracer.Start(ctx, h.OtelData.RequestNameOTEL)
	defer span.End()

	w.Header().Add("Content-Type", "application/json")

	var request dto.RequestWeather

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		http.Error(w, "invalid weather", http.StatusUnprocessableEntity)
		return
	}

	if len(request.City) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		http.Error(w, "invalid weather", http.StatusUnprocessableEntity)
		return
	}

	start := time.Now()
	response, err := services.GetWeatherApiService(request.City)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("response Weather:", time.Since(start))

	_, spanTime := h.OtelData.OTELTracer.Start(ctx, "SPAN_TIME_WEATHER: "+time.Since(start).String())
	spanTime.End()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
