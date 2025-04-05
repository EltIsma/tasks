package response

import (
	"encoding/json"
	"net/http"
)

func ResultErrJSON(res http.ResponseWriter, status int, body map[string]any) {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(status)

	_ = json.NewEncoder(res).Encode(body)
}

func ResultJSON(res http.ResponseWriter, status int, body []byte) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}
