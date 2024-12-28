package types

type Student struct {
	Id    int64  `json:"id"`
	Email string `json:"email" validate:"required"`
	Name  string `json:"name"  validate:"required"`
	Age   int    `json:"age"  validate:"required"`
}
