package metrics

import (
	"net/http"

	"github.com/gorilla/mux"
	prometheus "github.com/prometheus/client_golang/prometheus"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	DefaultPath     = "/metrics"
	DefaultHostPort = ":8019"
)

type HTTPServerConfig struct {
	MetricsPath     string
	MetricsHostPort string
	Handler         http.Handler
	Registerer      *prometheus.Registerer
}

type HTTPServer struct {
	server *http.Server
}

func NewDefaultHTTPServerConfig() *HTTPServerConfig {
	defaultConfig := &HTTPServerConfig{
		MetricsPath:     DefaultPath,
		MetricsHostPort: DefaultHostPort,
		Handler:         promhttp.Handler(),
		Registerer:      &prometheus.DefaultRegisterer,
	}
	return defaultConfig
}

func NewDefaultHTTPServer(config *HTTPServerConfig) (*HTTPServer, error) {
	if config == nil {
		config = NewDefaultHTTPServerConfig()
	}

	route := mux.NewRouter()
	route.NotFoundHandler = config.Handler

	srv := HTTPServer{
		server: &http.Server{
			Addr:              config.MetricsHostPort,
			Handler:           http.Handler(route),
			ReadHeaderTimeout: 0,
		},
	}

	route.Handle(config.MetricsPath, config.Handler)
	return &srv, nil
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop() error {
	return s.server.Close()
}
