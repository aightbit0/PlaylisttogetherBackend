package playlist

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func (p *PlaylistCredentials) createJWT(us string) string {
	sec := time.Now().Unix()
	expireOn := sec + 7200
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": us,
		"exp":  expireOn,
	})

	tokenString, err := token.SignedString(p.hmacSampleSecret)
	if err != nil {
		fmt.Printf("failed to generate JWT Token")
	}

	return tokenString
}

func (p *PlaylistCredentials) checkJWT(jwtToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.hmacSampleSecret, nil
	})
	if err != nil {
		return nil, err
	}
	autherror := errors.New("not authorised")
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, autherror
}
