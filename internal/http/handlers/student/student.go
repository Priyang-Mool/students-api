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
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student // Declare a variable to hold the student data

		// Decode the JSON request body into the student variable
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) { // Check if the request body is empty
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil { // Check for other decoding errors
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

		lastID, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		slog.Info("User Created Successfully!")

		// Respond with a success message and HTTP status 201 Created
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastID})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student!", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Fatal("Error during conversion from tring to int64")
			return
		}

		student, e := storage.GetStudentById(intId)
		if e != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(e))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)
	}
}

func GetAll(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")

		students, err := storage.GetAllStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJSON(w, http.StatusOK, students)
	}
}

func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Updating a student with", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			log.Fatal("error converting string to int")
			return
		}

		var student types.Student

		err = json.NewDecoder(r.Body).Decode(&student)

		if err == io.EOF {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		if err != nil {
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

		updatedStudent, err := storage.UpdateStudent(intId, student.Name, student.Email, student.Age)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Student Updated Successfully!")
		response.WriteJSON(w, http.StatusOK, updatedStudent)
	}
}
