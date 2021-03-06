package handler

import (
	"encoding/json"
	"log"
	"mflight-api/application"
	"mflight-api/domain"
	"net/http"
)

// SensorMetricsHandler is struct to get sensor metrics
type SensorMetricsHandler struct {
	metricsCollector application.MetricsCollector
}

type response []metrics

type metrics struct {
	Unixtime    int64   `json:"unixtime"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Illuminance int16   `json:"illuminance"`
}

// NewSensorMetricsHandler creates a new SensorMetricsHandler based on domain.Sensor
func NewSensorMetricsHandler(c application.MetricsCollector) *SensorMetricsHandler {
	return &SensorMetricsHandler{c}
}

// ServeHTTP implements http.Handler
func (h *SensorMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, err := h.metricsCollector.CollectTimeSeriesMetrics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := convert(m)

	bytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	successResponse(w, bytes)
}

func convert(ts domain.TimeSeriesMetrics) response {
	l := make([]metrics, len(ts))
	for i, m := range ts {
		l[i] = metrics{
			Unixtime:    m.Time.Unix(),
			Temperature: float32(m.Temperature),
			Humidity:    float32(m.Humidity),
			Illuminance: int16(m.Illuminance),
		}
	}
	return response(l)
}

func successResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(bytes)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}
