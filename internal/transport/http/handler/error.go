package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725732425-team-77001/zadanie-6105/internal/service"
)

func WriteValue(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	e.Encode(v)
}

func WriteReason(w http.ResponseWriter, code int, reason string) {
	WriteValue(w, code, map[string]string{
		"reason": reason,
	})
}

var ErrorCodes = map[service.ErrorType]int{
	service.ErrorTypeInvalid:      http.StatusBadRequest,
	service.ErrorTypeUnauthorized: http.StatusUnauthorized,
	service.ErrorTypeForbidden:    http.StatusForbidden,
	service.ErrorTypeNotExist:     http.StatusNotFound,
}

func HandleServiceError(w http.ResponseWriter, err error) {
	var serviceErr *service.Error
	if errors.As(err, &serviceErr) {
		code, ok := ErrorCodes[serviceErr.Type()]
		if ok {
			var reason strings.Builder
			reason.WriteString(serviceErr.Error())
			uerr := serviceErr.Unwrap()
			if uerr != nil {
				reason.WriteString(": ")
				reason.WriteString(uerr.Error())
			}
			WriteReason(w, code, reason.String())
			return
		}
		panic(fmt.Errorf("%s: %w", serviceErr.Error(), serviceErr.Unwrap()))
	}

	// Panic must be recovered by RecovererMiddleware.
	panic(err)
}
