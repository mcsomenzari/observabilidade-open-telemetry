package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"time"

	"github.com/mcsomenzari/temperature-system-by-cep-otel/internal/application/infra/web"
	"github.com/mcsomenzari/temperature-system-by-cep-otel/pkg/telemetry"

	"go.opentelemetry.io/otel"
)

// server B
func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := telemetry.InitProvider("serviceB", "otel-collector:4317")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown traceProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("microservice-tracer")
	otelData := &web.OtelData{
		RequestNameOTEL: "microservice-request-02",
		OTELTracer:      tracer,
	}
	server := web.NewServer(otelData)
	router := server.CreateServerB()

	go func() {
		log.Println("Starting server B on port", "8081")
		if err := http.ListenAndServe(":8081", router); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interruption.
	select {
	case <-sigCh:
		log.Printf("shutting down gracefully, CTRL+C pressed...")
		return
	case <-ctx.Done():
		log.Printf("shutting down due to other reason...")
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	defer func() {
		_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		shutdownCancel()
	}()
}
