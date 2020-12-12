package jsonfmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/appmux"
)

// Respond converts a Go value to JSON and sends it to the cli.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	// Set the status code for the request logger middlewarefunc.
	v := ctx.Value(appmux.KeyValues).(*appmux.Values)
	v.StatusCode = statusCode

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response value to JSON.
	res, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marshalling value to json")
	}

	// Respond with the provided JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(res); err != nil {
		return errors.Wrap(err, "writing to cli")
	}

	return nil
}

// RespondError sends an error reponse back to the cli.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {

	// If the error was of the type *ErrorRequest, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*appmux.ErrorRequest); ok {
		er := appmux.ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		fmt.Println("\tRespondError err is ", er)
		if err := Respond(ctx, w, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500.
	er := appmux.ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(ctx, w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
