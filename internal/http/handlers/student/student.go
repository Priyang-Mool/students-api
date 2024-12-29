package student

import (
	"encoding/json" // Package for JSON encoding and decoding
	"errors"        // Package for error handling
	"fmt"           // Package for formatted I/O
	"io"            // Package for I/O primitives
	"log"
	"log/slog" // Package for structured logging
	"net/http" // Package for HTTP client and server
	"strconv"

	"github.com/Priyang1310/Students-API-GO/internal/storage"
	"github.com/Priyang1310/Students-API-GO/internal/types"          // Importing custom types
	"github.com/Priyang1310/Students-API-GO/internal/utils/response" // Importing response utility functions
	"github.com/go-playground/validator/v10"                         // Importing the validator package for struct validation
)

// New returns an HTTP handler function for creating a new student
// This function handles the HTTP request to create a new student
// It validates the student data, creates a new student in the storage, and returns the created student's ID
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Declare a variable to hold the student data
		var student types.Student

		// Decode the JSON request body into the student variable
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) { // Check if the request body is empty
			// Return a bad request error if the request body is empty
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil { // Check for other decoding errors
			// Return a bad request error if there's a decoding error
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate the student struct using the validator package
		if err := validator.New().Struct(student); err != nil {
			// If validation fails, cast the error to ValidationErrors and respond with a validation error
			validateErr := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return // Ensure to return after sending the response
		}

		// Create a new student in the storage
		lastID, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			// Return an internal server error if there's an error creating the student
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		// Log a success message
		slog.Info("User Created Successfully!")

		// Respond with a success message and HTTP status 201 Created
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastID})
	}
}

// GetById returns an HTTP handler function for getting a student by ID
// This function handles the HTTP request to get a student by ID
// It retrieves the student from the storage and returns the student data
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the URL path
		id := r.PathValue("id")
		slog.Info("Getting a student!", slog.String("id", id))

		// Convert the ID to an integer
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			// Log a fatal error if there's an error converting the ID
			log.Fatal("Error during conversion from tring to int64")
			return
		}

		// Retrieve the student from the storage
		student, e := storage.GetStudentById(intId)
		if e != nil {
			// Return an internal server error if there's an error retrieving the student
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(e))
			return
		}

		// Respond with the student data
		response.WriteJSON(w, http.StatusOK, student)
	}
}

// GetAll returns an HTTP handler function for getting all students
// This function handles the HTTP request to get all students
// It retrieves all students from the storage and returns the student data
func GetAll(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log a message
		slog.Info("Getting all students")

		// Retrieve all students from the storage
		students, err := storage.GetAllStudents()
		if err != nil {
			// Return an internal server error if there's an error retrieving the students
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		// Respond with the student data
		response.WriteJSON(w, http.StatusOK, students)
	}
}

// Update returns an HTTP handler function for updating a student
// This function handles the HTTP request to update a student
// It validates the student data, updates the student in the storage, and returns the updated student data
func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the URL path
		id := r.PathValue("id")
		slog.Info("Updating a student with", slog.String("id", id))

		// Convert the ID to an integer
		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			// Log a fatal error if there's an error converting the ID
			log.Fatal("error converting string to int")
			return
		}

		// Declare a variable to hold the student data
		var student types.Student

		// Decode the JSON request body into the student variable
		err = json.NewDecoder(r.Body).Decode(&student)

		if err == io.EOF {
			// Return an internal server error if the request body is empty
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		if err != nil {
			// Return an internal server error if there's a decoding error
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Validate the student struct using the validator package
		if err := validator.New().Struct(student); err != nil {
			// If validation fails, cast the error to ValidationErrors and respond with a validation error
			validateErr := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return // Ensure to return after sending the response
		}

		// Update the student in the storage
		updatedStudent, err := storage.UpdateStudent(intId, student.Name, student.Email, student.Age)

		if err != nil {
			// Return an internal server error if there's an error updating the student
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Log a success message
		slog.Info("Student Updated Successfully!")

		// Respond with the updated student data
		response.WriteJSON(w, http.StatusOK, updatedStudent)
	}
}

// DeleteById returns an HTTP handler function for deleting a student by ID
// This function handles the HTTP request to delete a student by ID
// It deletes the student from the storage and returns a success message
func DeleteById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the URL path
		id := r.PathValue("id")

		// Convert the ID to an integer
		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			// Return an internal server error if there's an error converting the ID
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Delete the student from the storage
		err = storage.DeleteStudentById(intId)

		if err != nil {
			// Return an internal server error if there's an error deleting the student
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		}

		// Respond with a success message
		response.WriteJSON(w, http.StatusOK, "student deleted successfully")
	}
}

// DeleteAll returns an HTTP handler function for deleting all students
// This function handles the HTTP request to delete all students
// It deletes all students from the storage and returns a success message
func DeleteAll(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log a message
		slog.Info("Deleting All Students!")

		// Delete all students from the storage
		err := storage.DeleteAllStudents()
		if err != nil {
			// Return an internal server error if there's an error deleting the students
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Respond with a success message
		response.WriteJSON(w, http.StatusOK, "deleted all students successfully")
	}
}