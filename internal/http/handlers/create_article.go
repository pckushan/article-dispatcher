package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/models"
	"article-dispatcher/internal/domain/services"
	"article-dispatcher/internal/http/responses"

	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"

	"encoding/json"
	"fmt"
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
		ac.ErrorHandler.Handle(request.Context(), writer, InvalidPayload{
			fmt.Errorf("error decoding request body due to, %w", err)})
		return
	}

	// validate request struct
	if err = validate(&article); err != nil {
		ac.ErrorHandler.Handle(request.Context(), writer, ValidationError{
			fmt.Errorf("invalid request body due to, %w", err)})
		return
	}

	err = ac.ArticleService.Create(request.Context(), &article)
	if err != nil {
		ac.ErrorHandler.Handle(request.Context(), writer, fmt.Errorf("error marshaling response data, %w", err))
		return
	}

	resp := responses.SuccessResponse{}
	resp.Data.ID = article.Id
	r, err := json.Marshal(resp)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error marshaling response data due to, %s", err))
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_, err = writer.Write(r)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}

// validate - income request validator for the article payload
func validate(article *models.Article) error {
	return validator.New().Struct(article)
}
