package http

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/services"
	"article-dispatcher/internal/http/handlers"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gorilla/mux"

	"context"
	"fmt"
	"net/http"
)

type Router struct {
	server *http.Server
	Conf   *RouterConfig
	logger logger.Logger
}

func (r *Router) Init(l logger.Logger, articleService services.ArticleService, latencyReport *prometheus.SummaryVec) {
	muxRouter := mux.NewRouter()
	r.logger = l

	errorHandler := handlers.ErrorHandler{Log: l}

	r.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", r.Conf.Host),
		Handler:      muxRouter,
		ReadTimeout:  r.Conf.Timeouts.Read,
		WriteTimeout: r.Conf.Timeouts.Write,
		IdleTimeout:  r.Conf.Timeouts.Idle,
	}
	mw := handlers.Middleware{Logger: l}
	muxRouter.Use(mw.MiddleFunc)

	muxRouter.Handle(
		"/articles",
		handlers.ArticleCreateHandler{
			Log:                  l,
			ArticleService:       articleService,
			ErrorHandler:         errorHandler,
			RequestLatencyReport: latencyReport,
		}).Methods(http.MethodPost)

	muxRouter.Handle(
		"/articles/{id}",
		handlers.ArticleGetHandler{
			Log:                  l,
			ArticleService:       articleService,
			ErrorHandler:         errorHandler,
			RequestLatencyReport: latencyReport,
		}).Methods(http.MethodGet)
	muxRouter.Handle(
		"/tags/{tagName}/{date}",
		handlers.ArticleFilterHandler{
			Log:                  l,
			ArticleService:       articleService,
			ErrorHandler:         errorHandler,
			RequestLatencyReport: latencyReport,
		}).Methods(http.MethodGet)
}

func (r *Router) Start() error {
	r.logger.Info(fmt.Sprintf("server starting on port: %s", r.Conf.Host))
	if err := r.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (r *Router) Stop() error {
	c, fn := context.WithTimeout(context.Background(), r.Conf.Timeouts.ShoutDownWait)
	defer fn()
	r.logger.Info(fmt.Sprintf("server shutting down on port: %s", r.Conf.Host))
	return r.server.Shutdown(c)
}
