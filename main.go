package main

import (
	"fmt"
	"log"
	"mflight-exporter/config"
	"mflight-exporter/handler"
	"mflight-exporter/infrastructure/mflight"
	"mflight-exporter/infrastructure/prometheus/collector"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sensor := mflight.NewMfLightSensor(c.MfLight.URL, c.MfLight.MobileID)

	h := handler.NewSensorMetricsHandler(sensor)
	http.Handle("/getSensorMetrics", h)

	col := collector.NewMfLightCollector(sensor)
	prometheus.MustRegister(col)
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
}
