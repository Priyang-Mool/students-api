package response

import (
	"encoding/json" // Package for encoding and decoding JSON
	"fmt"           // Package for formatted I/O
	"net/http"      // Package for HTTP client and server
	"strings"       // Package for string manipulation

	"github.com/go-playground/validator/v10" // Package for data validation
)

// Response struct defines the structure of the JSON response
// It contains two fields: Status and Error
type Response struct {
	Status string // Status of the response (e.g., "Error" or "Success")
	Error  string // Error message if applicable
}

// WriteJSON writes a JSON response to the http.ResponseWriter
// It takes three parameters: http.ResponseWriter, HTTP status code, and data to be written
// It sets the Content-Type header to application/json, writes the HTTP status code, and encodes the data as JSON
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Write the HTTP status code
	w.WriteHeader(status)

	// Encode the data as JSON and write it to the response
	return json.NewEncoder(w).Encode(data)
}

// GeneralError creates a generic error response
// It takes an error as input and returns a Response struct with the error message
func GeneralError(err error) Response {
	// Create a Response struct with the error message
	return Response{
		Status: "Error", // Set status to "Error"
		Error:  err.Error(), // Include the error message
	}
}

// ValidationError creates a validation error response based on validation errors
// It takes a slice of validation errors as input and returns a Response struct with the error messages
func ValidationError(errs validator.ValidationErrors) Response {
	// Create a slice to hold individual error messages
	var errMsgs []string

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