package dtos

type LoginDTO struct {
	Email			string		`json:"email"`
	Password	string		`json:"password"`
}

type LoginTokenDTO struct {
	Token			string		`json:"jwt"`
}
