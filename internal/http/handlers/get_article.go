package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/services"

	"github.com/gorilla/mux"

	"encoding/json"
	"fmt"
	"net/http"
)

type ArticleGetHandler struct {
	Log            logger.Logger
	ArticleService services.ArticleService
	ErrorHandler   ErrorHandler
}

func (ag ArticleGetHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// capture query params
	vars := mux.Vars(request)
	articleID := vars[PathParameterArticleID]

	article, err := ag.ArticleService.Get(request.Context(), articleID)
	if err != nil {
		ag.Log.Error(fmt.Sprintf("error fetching article data due to, %s", err))
		ag.ErrorHandler.Handle(request.Context(), writer, err)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	r, err := json.Marshal(article)
	if err != nil {
		ag.Log.Error(fmt.Sprintf("error marshaling response data due to, %s", err))
	}
	_, err = writer.Write(r)
	if err != nil {
		ag.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}
