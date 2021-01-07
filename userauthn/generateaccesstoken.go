package userauthn

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn/userauthnvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/usermgmt/uservalue"
)

const (
	secretKey     = "kubernetssecret"
	tokenDuration = 1 * time.Minute
)

type customClaims struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

func generateAccessToken(ctx context.Context, db *sqlx.DB, ctreq *userauthnvalue.CreateTokenRequest) (*userauthnvalue.CreateTokenResponse, error) {
	vyu := ctreq.VerifyUser
	uFound, err := uservalue.FindByEmail(ctx, db, vyu.Email)
	log.Printf("user found by email is %+v", uFound)
	if err != nil {
		return nil, muxhttp.NewRequestError(
			errors.Wrap(err, fmt.Sprintf("NotFound Error: User not found for email %v", vyu.Email)),
			http.StatusNotFound)
	}
	log.Printf("we were actually able to find the user email %v\n", uFound.Email)
	err = PasswordMatch(uFound.PasswordHash, vyu.Password)
	log.Printf("Password match Err is %v\n", err)
	if err != nil {
		return nil, err
	}

	claims := customClaims{
		Email: uFound.Email,
		Roles: uFound.Roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
			Issuer:    "github.com/rsachdeva/illuminatingdeposits-rest",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, muxhttp.NewRequestError(
			errors.Wrap(err, fmt.Sprintf("Cannot generate access token for %v", vyu.Email)),
			http.StatusInternalServerError)
	}

	fmt.Println("signedToken generated finally is", signedToken)
	uaresp := userauthnvalue.CreateTokenResponse{
		VerifiedUser: &userauthnvalue.VerifiedUser{
			AccessToken: signedToken,
		},
	}
	return &uaresp, nil
}
