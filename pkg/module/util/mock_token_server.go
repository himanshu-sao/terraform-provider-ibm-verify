package util

import (
	"io"
	"log"
	"net/http"
)

// MockTokenHandler returns a handler function for the token endpoint.
func MockTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("MockTokenHandler invoked")

		// Log the incoming request details
		//log.Printf("Request Method: %s, URL: %s", r.Method, r.URL.String())
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %s. %s", err, string(body))
		} /*else {
			log.Printf("Request Body: %s", string(body))
		}*/

		// Set response headers and write the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"access_token": "mockAccessToken",
			"grant_id": "mockGrantID",
			"token_type": "Bearer",
			"expires_in": 3600
		}`
		//log.Printf("Response Body: %s", response)
		w.Write([]byte(response))
	}
}
