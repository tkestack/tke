package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

var parser = &jwt.Parser{}

func Parse(tokenStr string) (jwt.MapClaims, error) {
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}  else {
		return nil, fmt.Errorf("convert to MapClaims failed")
	}
}
