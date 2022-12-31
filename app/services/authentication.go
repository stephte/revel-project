package services

import (
	"revel-project/app/services/dtos"
	"github.com/google/uuid"
	"encoding/base64"
	"encoding/json"
	"crypto/sha256"
	"crypto/hmac"
	"strings"
	"errors"
	"fmt"
	"os"
)

type AuthService struct {
	*BaseService
}


func(this AuthService) GenerateJWT(header dtos.JWTHeaderDTO, payload dtos.JWTPayloadDTO) (string, dtos.ErrorDTO) {
	if !header.Exists() || !payload.Exists() {
		this.log.Errorf("Invalid JWT Params!\nHeaderDTO: %v;\nPayloadDTO: %v", header, payload)
		return "", dtos.CreateErrorDTO(errors.New("Error!"), 500, false)
	}
	
	payloadJSON, pMarshalErr := json.Marshal(payload)
	if pMarshalErr != nil {
		this.log.Error(pMarshalErr.Error())
		return "", dtos.CreateErrorDTO(errors.New("Error!"), 500, false)
	}

	headerJSON, hMarshalErr := json.Marshal(header)

	if hMarshalErr != nil {
		this.log.Error(hMarshalErr.Error())
		return "", dtos.CreateErrorDTO(errors.New("Error!"), 500, false)
	}

	encodedHeader := base64.StdEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.StdEncoding.EncodeToString(payloadJSON)

	signature := this.generateSignature(encodedHeader, encodedPayload)

	token := fmt.Sprintf("%s.%s.%s", encodedHeader, encodedPayload, signature)

	return token, dtos.ErrorDTO{}
}


func(this *AuthService) ValidateJWT(jwt string, isPWReset bool) (bool, dtos.ErrorDTO) {
	splitJWT := strings.Split(jwt, ".")

	if len(splitJWT) != 3 {
		this.log.Warnf("Token not in 3 parts: %s", jwt)
		return this.invalidTokenErr()
	}

	encodedHeader := splitJWT[0]
	encodedPayload := splitJWT[1]
	jwtSignature := splitJWT[2]

	payloadJSON, decodeErr := base64.StdEncoding.DecodeString(encodedPayload)

	if decodeErr != nil {
		this.log.Error(decodeErr.Error())
		return this.invalidTokenErr()
	}

	var payload dtos.JWTPayloadDTO
	marshalErr := json.Unmarshal(payloadJSON, &payload)
	if marshalErr != nil {
		this.log.Error(marshalErr.Error())
		return this.invalidTokenErr()
	}

	userKey, parseErr := uuid.Parse(payload.Key)
	if parseErr != nil {
		this.log.Errorf("Error parsing UUID: %s", payload.Key)
		return this.invalidTokenErr()
	}

	findErr := this.setCurrentUser(userKey)
	if findErr != nil {
		this.log.Errorf("User not found: %s", payload.Key)
		return this.invalidTokenErr()
	}

	signature := this.generateSignature(encodedHeader, encodedPayload)

	if jwtSignature != signature {
		this.log.Errorf("Signatures do not match: %s::%s", jwtSignature, signature)
		return this.invalidTokenErr()
	}

	if isPWReset != payload.PRT {
		this.log.Error("Resets don't match")
		return this.invalidTokenErr()
	} else if !payload.IsActive() {
		this.log.Error("Token expired")
		return false, dtos.ErrorDTO{}
	}

	return true, dtos.ErrorDTO{}
}


func(this *AuthService) generateSignature(header string, payload string) string {
	strToSign := fmt.Sprintf("%s%s", header, payload)

	signingKey := fmt.Sprintf("%s%s", os.Getenv("AUTH_SECRET"), this.currentUser.EncryptedPassword)

	hasher := hmac.New(sha256.New, []byte(signingKey))
	hasher.Write([]byte(strToSign))

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}


func(this AuthService) invalidTokenErr() (bool, dtos.ErrorDTO) {
	return false, dtos.AccessDeniedError()
}
