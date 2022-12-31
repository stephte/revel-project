package tests

import (
	"revel-project/app/services/dtos"
	"encoding/json"
	"strings"
	"bytes"
)

type LoginTest struct {
	BaseTest
}

func(t *LoginTest) Before() {
	t.createTestAuthUsers()
}

func (t *LoginTest) After() {
	t.CleanupAuth()
}

func(this *LoginTest) TestValidLogin() {
	regularData, jsonErr := json.Marshal(
		map[string]interface{}{
		"email": "regular@test.com",
		"password": "testpassword9",
	})
	if jsonErr != nil {
		panic(jsonErr)
	}

	tokenDTO, errDTO := this.handleLogin(regularData)

	this.Assert(!errDTO.Exists())
	this.Assert(tokenDTO.Token != "")

	splitJWT := strings.Split(tokenDTO.Token, ".")
	this.Assert(len(splitJWT) == 3)
}


func(this *LoginTest) TestInvalidLogin() {
	regularData, jsonErr := json.Marshal(
		map[string]interface{}{
		"email": "fakeemail1@test.com",
		"password": "password",
	})
	if jsonErr != nil {
		panic(jsonErr)
	}

	tokenDTO, errDTO := this.handleLogin(regularData)

	this.Assert(errDTO.Exists())
	this.Assert(tokenDTO.Token == "")
}


func(this *LoginTest) handleLogin(data []byte) (dtos.LoginTokenDTO, dtos.ErrorDTO) {
	resp, reqErr := this.Client.Post("http://localhost:9000/auth/login", "application/json", bytes.NewBuffer(data))
	if reqErr != nil {
		panic(reqErr)
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	
	bytes := buf.Bytes()

	var tokenData dtos.LoginTokenDTO
	var errData dtos.ErrorDTO

	jsonErr := json.Unmarshal(bytes, &tokenData)
	if jsonErr != nil {
		panic(jsonErr)
	}

	jsonErr = json.Unmarshal(bytes, &errData)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return tokenData, errData
}
