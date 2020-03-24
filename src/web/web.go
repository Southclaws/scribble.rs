package web

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// StatusNotFound writes a pretty error
func StatusNotFound(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	errToWriter(w, err)
}

// StatusInternalServerError writes a pretty error and logs
func StatusInternalServerError(w http.ResponseWriter, err error) {
	zap.L().Error("bad error", zap.Error(err))
	w.WriteHeader(http.StatusInternalServerError)
	errToWriter(w, err)
}

// StatusUnauthorized writes a pretty error
func StatusUnauthorized(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	errToWriter(w, err)
}

// StatusNotAcceptable writes a pretty error
func StatusNotAcceptable(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotAcceptable)
	errToWriter(w, err)
}

// StatusBadRequest writes a pretty error
func StatusBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotAcceptable)
	errToWriter(w, err)
}

func errToWriter(w io.Writer, err error) {
	if e := json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}{
		Message: err.Error(),
		Error:   err,
	}); e != nil {
		zap.L().Fatal(
			"failed to write error response",
			zap.Error(e),
			zap.Any("original_error", err))
	}
}
