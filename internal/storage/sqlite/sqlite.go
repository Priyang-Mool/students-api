package sqlite

import (
	"database/sql" // Import the database/sql package for SQL database operations

	"github.com/Priyang1310/Students-API-GO/internal/config" // Import the config package for application configuration
	_ "github.com/mattn/go-sqlite3"                          // Import the SQLite driver for database operations
)

// Sqlite struct represents a SQLite database connection
type Sqlite struct {
	Db *sql.DB // Db is a pointer to the sql.DB type, which represents a database connection
}

// New function initializes a new Sqlite instance
// It takes a configuration object as an argument and returns a pointer to Sqlite and an error
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

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?) ")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
