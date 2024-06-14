package web

import "go.opentelemetry.io/otel/trace"

type WebServer struct {
	OtelData *OtelData
}

type OtelData struct {
	RequestNameOTEL string
	OTELTracer      trace.Tracer
}

func NewServer(otelData *OtelData) *WebServer {
	return &WebServer{
		OtelData: otelData,
	}
}
