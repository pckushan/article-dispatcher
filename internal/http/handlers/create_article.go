package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/models"
	"article-dispatcher/internal/domain/services"
	"article-dispatcher/internal/http/responses"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

type ArticleCreateHandler struct {
	Log                  logger.Logger
	ArticleService       services.ArticleService
	ErrorHandler         ErrorHandler
	RequestLatencyReport *prometheus.SummaryVec
}

func (ac ArticleCreateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	var err error
	defer func() {
		ac.RequestLatencyReport.
			With(map[string]string{"endpoint": "add_article", "error": fmt.Sprintf(`%t`, err != nil)}).
			Observe(float64(time.Since(start).Microseconds()))
	}()
	var article models.Article
	err = json.NewDecoder(request.Body).Decode(&article)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error decoding request body due to, %s", err))
		errN := InvalidPayload{err}
		ac.ErrorHandler.Handle(request.Context(), writer, errN)
		return
	}

	err = ac.ArticleService.Create(request.Context(), &article)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error marshaling response data due to, %s", err))
		ac.ErrorHandler.Handle(request.Context(), writer, err)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	resp := responses.SuccessResponse{}
	resp.Data.ID = article.Id
	r, err := json.Marshal(resp)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error marshaling response data due to, %s", err))
	}
	_, err = writer.Write(r)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}
