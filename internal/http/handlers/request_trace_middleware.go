package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"context"
	"fmt"
	"github.com/google/uuid"

	"net/http"
)

type Middleware struct {
	Logger logger.Logger
}

func (mw Middleware) MiddleFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		traceID := uuid.New()
		request = request.WithContext(context.WithValue(request.Context(), ParamTraceID, traceID))
		mw.Logger.Debug(fmt.Sprintf("request received with trace-id:[%s] url [%s]", traceID, request.URL))
		handler.ServeHTTP(writer, request)
	})
}
