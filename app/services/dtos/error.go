package dtos



type ErrorResponse struct {
	Err 			string		`json:"error"`
	Relogin		bool			`json:"relogin"` // use this to tell the UI to re-authenticate (if login time has expired)...?
}