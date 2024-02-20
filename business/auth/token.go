package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken() error {
	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "12345678",
			Issuer:    "simplebank",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)), // i year
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate private key error: %w", err)
	}

	method := jwt.GetSigningMethod(jwt.SigningMethodRS256.Name)
	token := jwt.NewWithClaims(method, claims)
	str, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("sign token error: %w", err)
	}

	fmt.Println("### TOKEN ###")
	fmt.Println(str)
	fmt.Print("\n")

	return nil
}
