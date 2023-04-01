package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ArticleFilterHandler struct {
	Log            logger.Logger
	ArticleService services.ArticleService
	ErrorHandler   ErrorHandler
}

// ServeHTTP return a success response with the tagged article payload,
// if errors occur it will be sent to the error handler
func (af ArticleFilterHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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
		af.Log.Error(fmt.Sprintf("error marshalling response data due to, %s", err))
	}
	_, err = writer.Write(r)
	if err != nil {
		af.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}
