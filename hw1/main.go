package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/simonnik/GB_observability/hw1/middleware"
)

func main() {
	r := mux.NewRouter()
	web := http.Server{
		Handler: r,
		Addr:    ":8080",
	}
	metricsMiddleware := middleware.NewMetricsMiddleware()

	r.Use(metricsMiddleware.Metrics)

	r.HandleFunc("/identity", GetIdentityHandler).Methods(http.MethodPost)
	r.HandleFunc("/alert", AlertHandler).Methods(http.MethodGet)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != http.ErrServerClosed {
			panic(fmt.Errorf("error on listen and serve: %v", err))
		}
	}()
	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}

// GetIdentityHandler ...
func GetIdentityHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("token") == "admin_secret_token" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func AlertHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert"))
}
