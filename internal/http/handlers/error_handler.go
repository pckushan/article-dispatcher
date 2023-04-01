package handlers

import (
	"article-dispatcher/internal/adaptors/cache"
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/http/responses"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type ErrorHandler struct {
	Log logger.Logger
}

type internalErrorFields struct {
	code           int
	httpStatusCode int
	trace          string
}

func (e *ErrorHandler) Handle(ctx context.Context, writer http.ResponseWriter, err error) {
	e.Log.Error(fmt.Sprintf("error executing request due to : %s", err.Error()))
	errorBody := e.createErrorResponse(ctx, err)
	bodyByt, err := json.Marshal(errorBody)
	if err != nil {
		e.Log.Error(fmt.Sprintf(`failed to encode error response with, %s`, err))
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(errorBody.StatusCode)
	_, err = writer.Write(bodyByt)
	if err != nil {
		e.Log.Error(fmt.Sprintf(`failed to write error response due to, %s`, err))
	}
}

func (e *ErrorHandler) createErrorResponse(ctx context.Context, err error) responses.ErrorResponse {
	errorFields := mapError(err)
	return responses.ErrorResponse{
		StatusCode:  errorFields.httpStatusCode,
		Code:        errorFields.code,
		Description: errorFields.trace,
		Trace:       ctx.Value(ParamTraceID).(uuid.UUID).String(),
	}
}

func mapError(err error) internalErrorFields {
	switch err.(type) {
	case cache.TypeAssertError:
		return internalErrorFields{
			code:           CacheTypeAssertError,
			httpStatusCode: http.StatusBadRequest,
			trace:          err.Error(),
		}
	case cache.RepositoryError:
		return internalErrorFields{
			code:           RepositoryError,
			httpStatusCode: http.StatusBadRequest,
			trace:          err.Error(),
		}
	case cache.InvalidDataError:
		return internalErrorFields{
			code:           InvalidRequestDataError,
			httpStatusCode: http.StatusBadRequest,
			trace:          err.Error(),
		}
	default:
		if errors.Unwrap(err) == nil {
			return internalErrorFields{
				code:           UnknownError,
				httpStatusCode: http.StatusBadRequest,
				trace:          "something went wrong.",
			}
		}
		return mapError(errors.Unwrap(err)) // todo ??
	}
}
