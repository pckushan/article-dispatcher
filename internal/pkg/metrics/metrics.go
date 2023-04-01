package metrics

import (
	"article-dispatcher/internal/domain/adaptors/logger"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"context"
	"fmt"
	"net/http"
)

type RouterMetrics struct {
	server *http.Server
	Conf   *MetricConfig
	Logger logger.Logger
}

var RequestLatency *prometheus.SummaryVec

// InitMetrics init server and metrics reports
func (rm *RouterMetrics) InitMetrics() error {
	muxRouter := mux.NewRouter()
	muxRouter.Handle(`/metrics`, promhttp.Handler())
	rm.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", rm.Conf.HTTP.Host),
		Handler: muxRouter,
	}

	RequestLatency = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: rm.Conf.System,
		Subsystem: rm.Conf.SubSystem,
		Name:      "request_latency_micro",
		Help:      "http_request_latency",
	}, []string{"endpoint", "error"})

	prometheus.MustRegister(RequestLatency)

	return nil
}

func (rm *RouterMetrics) Start() error {
	rm.Logger.Info(fmt.Sprintf(`metrics server starting on port :%s`, rm.Conf.HTTP.Host))
	if err := rm.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (rm *RouterMetrics) Stop() error {
	c, fn := context.WithTimeout(context.Background(), rm.Conf.HTTP.ShutdownWait)
	defer fn()
	rm.Logger.Info(fmt.Sprintf(`metrics server shutting down on port %s`, rm.Conf.HTTP.Host))
	return rm.server.Shutdown(c)
}
