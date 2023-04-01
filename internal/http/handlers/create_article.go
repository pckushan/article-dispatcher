package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/models"
	"article-dispatcher/internal/domain/services"
	"article-dispatcher/internal/http/responses"
	"encoding/json"
	"fmt"
	"net/http"
)

type ArticleCreateHandler struct {
	Log            logger.Logger
	ArticleService services.ArticleService
	ErrorHandler   ErrorHandler
}

func (ac ArticleCreateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var article models.Article
	err := json.NewDecoder(request.Body).Decode(&article)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error decoding request body due to, %s", err))
	}

	err = ac.ArticleService.Create(request.Context(), article)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error marshalling response data due to, %s", err))
		ac.ErrorHandler.Handle(request.Context(), writer, err)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	resp := responses.SuccessResponse{}
	resp.Data.ID = article.Id
	r, err := json.Marshal(resp)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error marshalling response data due to, %s", err))
	}
	_, err = writer.Write(r)
	if err != nil {
		ac.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}
