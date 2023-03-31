package http

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/http/handlers"
	"context"
	"fmt"

	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	server *http.Server
	Conf   *RouterConfig
	logger logger.Logger
}

func (r *Router) Init(l logger.Logger) {
	muxRouter := mux.NewRouter()
	r.logger = l

	r.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", r.Conf.Host),
		Handler:      muxRouter,
		ReadTimeout:  r.Conf.Timeouts.Read,
		WriteTimeout: r.Conf.Timeouts.Write,
		IdleTimeout:  r.Conf.Timeouts.Idle,
	}
	mw := handlers.Middleware{Logger: l}
	muxRouter.Use(mw.MiddleFunc)

	muxRouter.Handle("/articles", handlers.ArticleCreateHandler{Log: l}).Methods(http.MethodPost)
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
