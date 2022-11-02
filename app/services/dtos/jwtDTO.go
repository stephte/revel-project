package dtos

type JWTPayloadDTO struct {
	Key					string			`json:"key"`
	Expiration	int64				`json:"exp"`
	Issuer			string			`json:"iss"`
}

type JWTHeaderDTO struct {
	Algorithm		string			`json:"alg"`
	Type				string			`json:"typ"`
}

func(dto JWTPayloadDTO) Exists() bool {
	return dto.Key != "" && dto.Expiration != 0 && dto.Issuer != ""
}

func(dto JWTHeaderDTO) Exists() bool {
	return dto.Algorithm != "" && dto.Type != ""
}
