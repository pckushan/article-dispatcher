package handlers

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/domain/models"
	"article-dispatcher/internal/http/responses"
	"encoding/json"
	"fmt"
	"net/http"
)

type ArticleCreateHandler struct {
	Log logger.Logger
}

func (ar ArticleCreateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var a models.Article
	_ = json.NewDecoder(request.Body).Decode(&a)

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	resp := responses.SuccessResponse{}
	resp.Data.ID = a.Id
	r, err := json.Marshal(resp)
	if err != nil {
		ar.Log.Error(fmt.Sprintf("error marshalling response data due to, %s", err))
	}
	_, err = writer.Write(r)
	if err != nil {
		ar.Log.Error(fmt.Sprintf("error writing to response due to, %s", err))
	}
}
