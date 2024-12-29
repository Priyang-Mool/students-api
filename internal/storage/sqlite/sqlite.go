package sqlite

import (
	"database/sql" // Import the database/sql package for SQL database operations
	"fmt"
	"log/slog"

	"github.com/Priyang1310/Students-API-GO/internal/config" // Import the config package for application configuration
	"github.com/Priyang1310/Students-API-GO/internal/types"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver for database operations
)

// Sqlite struct represents a SQLite database connection
type Sqlite struct {
	Db *sql.DB // Db is a pointer to the sql.DB type, which represents a database connection
}

// New function initializes a new Sqlite instance
// It takes a configuration object as an argument and returns a pointer to Sqlite and an error
// This function is used to establish a connection to the SQLite database
func New(cfg *config.Config) (*Sqlite, error) {
	// Open a new database connection using the SQLite driver and the storage path from the config
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		// If there is an error opening the database, return nil and the error
		return nil, err
	}

	// Execute a SQL command to create the 'students' table if it does not already exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT, 
		email TEXT, 
		age INTEGER 
	)`)

	if err != nil {
		// If there is an error executing the SQL command, return nil and the error
		return nil, err
	}

	// Return a new Sqlite instance with the established database connection
	return &Sqlite{
		Db: db, // Assign the database connection to the Db field of the Sqlite struct
	}, nil
}

// CreateStudent function creates a new student in the database
// It takes the student's name, email, and age as arguments and returns the ID of the newly created student and an error
// This function is used to insert a new student into the 'students' table
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	// Prepare a SQL statement to insert a new student into the 'students' table
	stmt, err := s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?) ")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	// Execute the prepared SQL statement with the provided student data
	result, err := stmt.Exec(name, email, age)

	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created student
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetStudentById function retrieves a student from the database by their ID
// It takes the student's ID as an argument and returns the student data and an error
// This function is used to select a student from the 'students' table by their ID
func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	// Prepare a SQL statement to select a student from the 'students' table by their ID
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	// Execute the prepared SQL statement with the provided student ID
	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {
		// If the student is not found, return an error
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student not found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, err
	}

	return student, nil
}

// GetAllStudents function retrieves all students from the database
// It returns a slice of student data and an error
// This function is used to select all students from the 'students' table
func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
	// Prepare a SQL statement to select all students from the 'students' table
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students")
	slog.Info("Get all students method called")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared SQL statement
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Initialize a slice to store the student data
	var students []types.Student

	// Iterate over the rows and scan the student data
	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		// Append the student data to the slice
		students = append(students, student)
	}

	return students, nil
}

// UpdateStudent function updates a student in the database
// It takes the student's ID, name, email, and age as arguments and returns the updated student data and an error
// This function is used to update a student in the 'students' table
func (s *Sqlite) UpdateStudent(id int64, name string, email string, age int) (types.Student, error) {
	// Prepare a SQL statement to update a student in the 'students' table
	slog.Info("Updating a student")
	stmt, err := s.Db.Prepare("UPDATE students SET name=?, email=?,age=? WHERE id=?")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	// Initialize a student struct to store the updated data
	var student types.Student

	// Execute the prepared SQL statement with the provided student data
	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return types.Student{}, err
	}

	// Get the number of rows affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.Student{}, err
	}
	// If no rows were affected, return an error
	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("no rows affected")
	}

	// Update the student struct with the provided data
	student.Age = age
	student.Email = email
	student.Name = name
	student.Id = id

	return student, nil
}

// DeleteStudentById function deletes a student from the database by their ID
// It takes the student's ID as an argument and returns an error
// This function is used to delete a student from the 'students' table by their ID
func (s *Sqlite) DeleteStudentById(id int64) error {
	// Prepare a SQL statement to delete a student from the 'students' table by their ID
	slog.Info("Deleting a student")
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id=?")
	if err != nil {
		return err
	}

	// Execute the prepared SQL statement with the provided student ID
	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

// DeleteAllStudents function deletes all students from the database
// It returns an error
// This function is used to delete all students from the 'students' table
func (s *Sqlite) DeleteAllStudents() error {
	// Prepare a SQL statement to delete all students from the 'students' table
	stmt, err := s.Db.Prepare("DELETE FROM students")

	if err != nil {
		return err
	}

	// Execute the prepared SQL statement
	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	return nil
}