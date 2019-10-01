package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/pabluk/promql-share/internal/application"
	"github.com/pabluk/promql-share/internal/diagnostics"
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	logger.Infof(
		"Starting the application: v.%s, commit:%s, buildTime:%s",
		diagnostics.Version, diagnostics.Commit, diagnostics.BuildTime,
	)

	port := os.Getenv("PORT")
	if port == "" {
		logger.Fatal("Port is not provided")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", application.HomeHandler(logger))

	prometheus, err := url.Parse("http://localhost:9090/")
	proxy := httputil.NewSingleHostReverseProxy(prometheus)
	r.PathPrefix("/static/").Handler(application.PrometheusHandler(logger, proxy)).Methods("GET")

	r.HandleFunc("/share", application.NewShareHandler(logger)).Methods("POST")
	r.HandleFunc("/share/{shareID:[0-9a-zA-Z]+}", application.GetShareHandler(logger)).Methods("GET")

	r.HandleFunc("/healthz", diagnostics.LivenessHandler(logger))
	r.HandleFunc("/readyz", diagnostics.ReadinessHandler(logger))

	r.HandleFunc("/{shareID:[0-9a-zA-Z]+}", application.GoToShareHandler(logger)).Methods("GET")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	shutdown := make(chan error, 1)

	server := http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: r,
	}

	go func() {
		err := server.ListenAndServe()
		shutdown <- err
	}()
	logger.Infof("The service is ready to listen and serve")

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			logger.Info("Got SIGINT...")
		case syscall.SIGTERM:
			logger.Info("Got SIGTERM...")
		}
	case <-shutdown:
		logger.Info("Got an error...")
	}

	logger.Infof("The service is stopping...")
	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Infof("Got an error during service shutdown: %v", err)
	}
}
