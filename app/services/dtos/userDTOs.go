package dtos



type UserDTO struct {
	BaseDTO
	FirstName 	string		`json:"firstName"`
	LastName		string		`json:"lastName"`
	Email				string		`json:"email"`
	Role				int				`json:"role"`
}

type CreateUserDTO struct {
	UserDTO
	Password		string	`json:"password"`
}
