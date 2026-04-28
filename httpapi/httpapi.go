package httpapi

import (
	"EmptyClassroom/service"
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, response service.APIResponse) {
	body, err := json.Marshal(response)
	if err != nil {
		status = http.StatusInternalServerError
		body = []byte(`{"code":500,"msg":"marshal failed","data":null}`)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func MethodNotAllowed(w http.ResponseWriter) {
	WriteJSON(w, http.StatusMethodNotAllowed, service.APIResponse{
		Code: 405,
		Msg:  "method not allowed",
		Data: nil,
	})
}

func Unauthorized(w http.ResponseWriter, msg string) {
	WriteJSON(w, http.StatusUnauthorized, service.APIResponse{
		Code: 401,
		Msg:  msg,
		Data: nil,
	})
}
