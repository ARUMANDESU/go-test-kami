package api

import (
	"net/http"
)

// logError logs an error message with the API's logger.
func (api API) logError(err error) {
	api.log.Error(err.Error())
}

// errorResponse sends a JSON response with the error message to the client.
func (api API) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := api.writeJSON(w, status, env, nil)
	if err != nil {
		api.logError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse logs the error and sends a generic error message to the client.
// It should be used when the server encounters an unexpected condition.
func (api API) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	api.logError(err)

	message := "the server encountered a problem and could not process your request"
	api.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse sends a 404 Not Found response to the client.
func (api API) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	api.errorResponse(w, r, http.StatusNotFound, message)
}

// failedValidationResponse sends a 422 Unprocessable Entity response to the client.
func (api API) failedValidationResponse(w http.ResponseWriter, r *http.Request, err error) {
	api.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// conflictResponse sends a 409 Conflict response to the client.
func (api API) conflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to handle the request due to a conflict with the current state of the resource"
	api.errorResponse(w, r, http.StatusConflict, message)
}
