package types

type Student struct {
	Id int32 
	Name string `validate:"required"` 
	Email string `validate:"required"`
	Age int `validate:"required"`
}
