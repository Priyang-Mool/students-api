package storage

import "github.com/Priyang1310/Students-API-GO/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	UpdateStudent(id int64, name string, email string, age int) (types.Student, error)
	DeleteStudentById(id int64) error
	DeleteAllStudents() error
}
