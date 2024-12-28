package types

type Student struct {
	Id    int
	Email string `validate:"required"`
	Name  string `validate:"required"`
	Age   int    `validate:"required"`
}
