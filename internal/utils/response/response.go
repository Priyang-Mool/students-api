package response

import (
	"encoding/json" // Package for encoding and decoding JSON
	"fmt"           // Package for formatted I/O
	"net/http"      // Package for HTTP client and server
	"strings"       // Package for string manipulation

	"github.com/go-playground/validator/v10" // Package for data validation
)

// Response struct defines the structure of the JSON response
type Response struct {
	Status string // Status of the response (e.g., "Error" or "Success")
	Error  string // Error message if applicable
}

// WriteJSON writes a JSON response to the http.ResponseWriter
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Write the HTTP status code
	w.WriteHeader(status)

	// Encode the data as JSON and write it to the response
	return json.NewEncoder(w).Encode(data)
}

// GeneralError creates a generic error response
func GeneralError(err error) Response {
	return Response{
		Status: "Error", // Set status to "Error"
		Error:  err.Error(), // Include the error message
	}
}

// ValidationError creates a validation error response based on validation errors
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string // Slice to hold individual error messages

	// Iterate over each validation error
	for _, err := range errs {
		// Check the type of validation error
		switch err.ActualTag() {
		case "required":
			// If the field is required but not provided, add a specific message
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			// For other types of validation errors, add a general invalid message
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	// Join all error messages into a single string and return the response
	return Response{
		Status: "Error", // Set status to "Error"
		Error: strings.Join(errMsgs, ", "), // Combine all error messages
	}
}