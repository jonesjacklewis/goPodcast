package api

import (
	"encoding/json"
	"net/http"
)

func (app *Application) writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func (app *Application) writeError(w http.ResponseWriter, status int, message string) error {
	errorMessage := map[string]any{
		"error":   true,
		"message": message,
	}

	return app.writeJson(w, status, errorMessage)
}

func (app *Application) createSuccessMessage(message string) map[string]any {
	return map[string]any{
		"error":   false,
		"message": message,
	}
}
