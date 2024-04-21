package responsehelper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Setup5xx setup Internal Server error response.
func Setup5xx(w http.ResponseWriter, err error) {
	SetupJSONResponse(w, http.StatusInternalServerError, "error", fmt.Errorf("internal server error: %w", err).Error())
}

// SetupJSONResponse is setup default JSON-response.
func SetupJSONResponse(w http.ResponseWriter, statusCode int, headerName string, headerBody string) {
	w.WriteHeader(statusCode)
	setupContentTypeJSON(w)
	resp := make(map[string]string)
	resp[headerName] = headerBody
	jsonResp, merr := json.Marshal(resp)
	if merr != nil {
		log.Fatalf("an error occurred while marshaling JSON: %s", merr.Error())
	}
	w.Write(jsonResp)
}

func setupContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
