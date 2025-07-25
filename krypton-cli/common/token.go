package common

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v4"
)

func GetClaimFromToken(token, claim string) (string, error) {
	t, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		if value, ok := claims[claim]; ok {
			return fmt.Sprint(value), nil
		}
	}
	return "", fmt.Errorf("no such claim: %s", claim)
}
