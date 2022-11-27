package dtos


type UserDTO struct {
	BaseDTO
	FirstName 	string
	LastName		string
	Email				string
	Role				int
}

type CreateUserDTO struct {
	UserDTO
	Password		string
}
