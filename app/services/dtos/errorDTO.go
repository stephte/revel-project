package dtos

type ErrorDTO struct {
	Err 			string		`json:"error"`
	Status		int				`json:"-"`
	Relogin		bool			`json:"relogin"` // tells the UI to re-authenticate
}

func CreateErrorDTO(err error, status int, relogin bool) ErrorDTO {
	if status == 0 {
		status = 400
	}

	return ErrorDTO{err.Error(), status, relogin}
}

func (dto ErrorDTO) Exists() bool {
	return dto.Err != ""
}
