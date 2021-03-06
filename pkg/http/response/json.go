/*
Package response provides helpers and utils for working with HTTP response
*/
package response

import (
	"context"
	"encoding/json"
	"net/http"
)

// JSON returns data as json response
func JSON(ctx context.Context, w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	// If there is nothing to marshal then set status code and return.
	if payload == nil {
		_, err := w.Write([]byte("{}"))
		return err
	}

	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)
	encoder.SetIndent("", "")

	if err := encoder.Encode(payload); err != nil {
		return err
	}

	Flush(w)

	return nil
}

// MustJSON returns data as json response
// will panic if unable to marshal payload into JSON object
// uses JSON internally
func MustJSON(ctx context.Context, w http.ResponseWriter, statusCode int, payload interface{}) {
	if err := JSON(ctx, w, statusCode, payload); err != nil {
		panic(err)
	}
}
