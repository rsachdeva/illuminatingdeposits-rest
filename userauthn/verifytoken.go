package userauthn

import (
	"fmt"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// valid validates the authorization.
func valid(authorization []string) error {
	if len(authorization) < 1 {
		return errors.New("no authorization header")
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	claims, err := verify(token)
	if err != nil {
		return err
	}
	email := claims.Email
	if len(email) < 1 {
		return errors.New("invalid token without email")
	}
	return nil
}

// Verify verifies the access token string and return a user claim if the token is valid
func verify(accessToken string) (*customClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}

			return []byte(secretKey), nil
		},
	)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			log.Printf("interceptor verify got invalid token ve is: %v", ve)
			if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
				return nil, errors.New(fmt.Sprintf("Token is expired. Please recreate token: %v", err))
			}
		}
		log.Printf("interceptor verify got invalid token: %v", err)
		return nil, errors.New(fmt.Sprintf("Verify got invalid token %v", err))
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	fmt.Printf("\nIn Verify claims are %+v\n", claims)
	return claims, nil
}
