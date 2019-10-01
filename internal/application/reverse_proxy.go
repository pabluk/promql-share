package application

import (
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

func PrometheusHandler(logger *logrus.Logger, p *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.URL)
		w.Header().Set("X-Ben", "Rad")
		p.ServeHTTP(w, r)
	})
}
