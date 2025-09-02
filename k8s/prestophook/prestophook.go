package prestophook

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	// DefaultCancelDelay is the default delay before cancelling the context
	DefaultCancelDelay = 10 * time.Minute

	// DefaultPort is the default port for the preStop hook server
	DefaultPort = "8029"
)

type PreStopHookConfig struct {
	Port         string
	PreStopDelay time.Duration
}

type PreStopHook struct {
	server      *http.Server
	config      PreStopHookConfig
	preStopFunc func()
}

func New(config PreStopHookConfig, preStopFunc func()) *PreStopHook {
	if config.Port == "" {
		config.Port = DefaultPort
	}
	return &PreStopHook{
		config:      config,
		server:      &http.Server{Addr: ":" + config.Port},
		preStopFunc: preStopFunc,
	}
}

func (p *PreStopHook) Start(ctx context.Context) {
	http.HandleFunc("/preStop", p.preStopHandler)

	go func() {
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listenAndServe(): %s", err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), DefaultCancelDelay)
	defer cancel()

	if err := p.server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %+v", err)
	}

	log.Println("server exited properly")
}

func (p *PreStopHook) preStopHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received preStop signal, executing preStop actions...")
	p.preStopFunc()

	time.Sleep(p.config.PreStopDelay)

	log.Println("preStop actions completed")
	w.WriteHeader(http.StatusOK)
}
