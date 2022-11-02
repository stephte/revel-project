package services

import (
	"revel-project/app/services/dtos"
	"github.com/google/uuid"
	"encoding/base64"
	"encoding/json"
	"crypto/sha256"
	"crypto/hmac"
	"math/rand"
	"strings"
	"errors"
	"time"
	"fmt"
	"os"
)

func generateJWT(header dtos.JWTHeaderDTO, payload dtos.JWTPayloadDTO) (string, error) {
	if !header.Exists() || !payload.Exists() {
		return "", errors.New("Invalid JWT params")
	}
	
	payloadJSON, pMarshalErr := json.Marshal(payload)
	if pMarshalErr != nil {
		return "", pMarshalErr
	}

	headerJSON, hMarshalErr := json.Marshal(header)

	if hMarshalErr != nil {
		return "", hMarshalErr
	}

	encodedHeader := base64.StdEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.StdEncoding.EncodeToString(payloadJSON)

	signature := generateSignature(encodedHeader, encodedPayload)

	token := fmt.Sprintf("%s.%s.%s", encodedHeader, encodedPayload, signature)

	return token, nil
}

func ValidateJWT(jwt string) (bool, error) {
	splitJWT := strings.Split(jwt, ".")

	if len(splitJWT) != 3 {
		return false, errors.New("Invalid Token")
	}

	encodedHeader := splitJWT[0]
	encodedPayload := splitJWT[1]
	jwtSignature := splitJWT[2]

	signature := generateSignature(encodedHeader, encodedPayload)

	if jwtSignature != signature {
		return false, errors.New("Invalid Token (possibly tampered with)")
	}

	payloadJSON, decodeErr := base64.StdEncoding.DecodeString(encodedPayload)

	if decodeErr != nil {
		panic(decodeErr)
	}

	return checkTokenActive(payloadJSON), nil
}

func checkTokenActive(payloadJSON []byte) bool {
	var payload dtos.JWTPayloadDTO
	err := json.Unmarshal(payloadJSON, &payload)

	if err != nil {
		panic(err)
	}

	return time.Now().Unix() < payload.Expiration
}

func generateSignature(header string, payload string) string {
	strToSign := fmt.Sprintf("%s%s", header, payload)

	hasher := hmac.New(sha256.New, []byte(os.Getenv("REVEL_SECRET")))
	hasher.Write([]byte(strToSign))

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func GetJWTUserKey(jwt string) uuid.UUID {
	splitJWT := strings.Split(jwt, ".")

	payloadJSON, decodeErr := base64.StdEncoding.DecodeString(splitJWT[1])
	if decodeErr != nil {
		panic(decodeErr)
	}

	var payload dtos.JWTPayloadDTO
	jsonErr := json.Unmarshal(payloadJSON, &payload)
	if jsonErr != nil {
		panic(jsonErr)
	}

	key, parseErr := uuid.Parse(payload.Key)

	if parseErr != nil {
		panic(parseErr)
	}

	return key
}

// time params in milliseconds
func killSomeTime(min int, max int) {
	// add a random amount of time (for security purposes)
	rand.Seed(time.Now().UnixNano())

	amount := rand.Intn(max-min) + min

	time.Sleep(time.Duration(amount) * time.Millisecond)
}
