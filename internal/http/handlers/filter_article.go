package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/services"
	"github.com/prometheus/client_golang/prometheus"
	"time"

	"github.com/gorilla/mux"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ArticleFilterHandler struct {
	Log                  logger.Logger
	ArticleService       services.ArticleService
	ErrorHandler         ErrorHandler
	RequestLatencyReport *prometheus.SummaryVec
}

// ServeHTTP return a success response with the tagged article payload,
// if errors occur it will be sent to the error handler
func (af ArticleFilterHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	var err error
	defer func() {
		af.RequestLatencyReport.
			With(map[string]string{"endpoint": "filter_article", "error": fmt.Sprintf(`%t`, err != nil)}).
			Observe(float64(time.Since(start).Microseconds()))
	}()

	// capture query params
	vars := mux.Vars(request)
	articleTag := vars[PathParameterTag]
	articleDate := vars[PathParameterDate]
	// todo validate date format
	date, err := strconv.Atoi(articleDate)
	if err != nil {
		af.Log.Error(fmt.Sprintf("error converting input date due to, %s", err))
		af.ErrorHandler.Handle(request.Context(), writer, err)
		return
	}
	taggedArticles, err := af.ArticleService.Filter(request.Context(), articleTag, date)
	if err != nil {
		af.Log.Error(fmt.Sprintf("error fetching tagged articles data due to, %s", err))
		af.ErrorHandler.Handle(request.Context(), writer, err)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	r, err := json.Marshal(taggedArticles)
	if err != nil {
		af.Log.Error(fmt.Sprintf("error marshaling response data due to, %s", err))
	}
	_, err = writer.Write(r)
	if err != nil {
		af.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}
