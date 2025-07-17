package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	e "github.com/jdpolicano/govault/internal/server/errors"
)

const staticServerError = `{"code":500,"error":"Internal Server Error"}`

type Response struct {
	Code  int    `json:"code"`
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func JSONResponse(w http.ResponseWriter, res Response) {
	if w == nil {
		return
	}
	payload, err := json.Marshal(res)
	if err != nil {
		fallbackServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(payload)
}

func NewResponse(code int, data any, error error) Response {
	if error != nil {
		return Response{code, data, error.Error()}
	}
	return Response{code, data, ""}
}

func NewInvalidBodyError() Response {
	return NewClientError(e.InvalidRequestBody)
}

func NewServerError(reason error) Response {
	return NewResponse(500, nil, reason)
}

func NewClientError(reason error) Response {
	return NewResponse(400, nil, reason)
}

func NewServerSuccess(data any) Response {
	return NewResponse(200, data, nil)
}

func fallbackServerError(w http.ResponseWriter, e error) {
	fmt.Println(e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(staticServerError))
}
