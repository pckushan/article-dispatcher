package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/services"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ArticleGetHandler struct {
	Log                  logger.Logger
	ArticleService       services.ArticleService
	ErrorHandler         ErrorHandler
	RequestLatencyReport *prometheus.SummaryVec
}

func (ag ArticleGetHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	var err error
	defer func() {
		ag.RequestLatencyReport.
			With(map[string]string{"endpoint": "get_article", "error": fmt.Sprintf(`%t`, err != nil)}).
			Observe(float64(time.Since(start).Microseconds()))
	}()

	// capture query params
	vars := mux.Vars(request)
	articleID := vars[PathParameterArticleID]
	if !validateArticleID(articleID) {
		err = fmt.Errorf("invalid article id format")
		ag.ErrorHandler.Handle(request.Context(), writer, ValidationError{err})
		return
	}

	article, err := ag.ArticleService.Get(request.Context(), articleID)
	if err != nil {
		err = fmt.Errorf("error fetching article data due to, %w", err)
		ag.ErrorHandler.Handle(request.Context(), writer, err)
		return
	}

	r, err := json.Marshal(article)
	if err != nil {
		ag.ErrorHandler.Handle(request.Context(), writer,
			ResponseMarshalError{fmt.Errorf("error marshaling response data, %w", err)})
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(r)
	if err != nil {
		ag.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}

// validateArticleID - converts the id path parameter into a integer and check for errors,
// if it converts into an integer it validates as a valid article id
func validateArticleID(id string) bool {
	_, err := strconv.Atoi(id)

	return err == nil
}
