package vars

import (
	"encoding/json"
	"net/http"
)

type HTTPStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func HTTPResp(w http.ResponseWriter, resp []byte) {
	w.WriteHeader(200)
	w.Write(resp)
}

func HTTPSuccess(w http.ResponseWriter, msg string) {
	status := HTTPStatus{
		Code:    "success",
		Message: msg,
		Status:  "success",
	}
	out, _ := json.Marshal(status)
	w.Write(out)
}

func HTTPError(w http.ResponseWriter, error string, msg string, code int) {
	status := HTTPStatus{
		Code:    error,
		Message: msg,
		Status:  "error",
	}
	out, _ := json.Marshal(status)
	http.Error(w, string(out), code)
}
