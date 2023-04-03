package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/services"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

	// validate input date
	if !validatePathDate(articleDate) {
		err = fmt.Errorf("invalid article date format")
		af.ErrorHandler.Handle(request.Context(), writer, ValidationError{err})
		return
	}

	// convert date into integer format `yyyymmdd`
	date, err := strconv.Atoi(articleDate)
	if err != nil {
		err = fmt.Errorf("error converting input date due to, %w", err)
		af.ErrorHandler.Handle(request.Context(), writer, err)
		return
	}
	taggedArticles, err := af.ArticleService.Filter(request.Context(), articleTag, date)
	if err != nil {
		err = fmt.Errorf("error fetching tagged articles data due to, %w", err)
		af.ErrorHandler.Handle(request.Context(), writer, err)
		return
	}

	r, err := json.Marshal(taggedArticles)
	if err != nil {
		af.ErrorHandler.Handle(request.Context(), writer,
			ResponseMarshalError{fmt.Errorf("error marshaling response data, %w", err)})
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(r)
	if err != nil {
		af.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}

func validatePathDate(date string) bool {
	_, err := time.Parse("20060102", date)

	return err == nil
}
