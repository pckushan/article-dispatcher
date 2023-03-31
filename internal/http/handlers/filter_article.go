package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type ArticleFilterHandler struct {
	Log            logger.Logger
	ArticleService services.ArticleService
}

func (af ArticleFilterHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// capture query params
	vars := mux.Vars(request)
	articleTag := vars[PathParameterTag]
	articleDate := vars[PathParameterDate]

	taggedArticles, err := af.ArticleService.Filter(request.Context(), articleTag, articleDate)
	if err != nil {
		// fixme handle error
		af.Log.Error(fmt.Sprintf("error fetching tagged articles data due to, %s", err))
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	r, err := json.Marshal(taggedArticles)
	if err != nil {
		af.Log.Error(fmt.Sprintf("error marshalling response data due to, %s", err))
	}
	_, err = writer.Write(r)
	if err != nil {
		af.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}
