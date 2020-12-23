package userauthnvalue

type VerifyUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifiedUser struct {
	AccessToken string `json:"access_token"`
}

type CreateTokenRequest struct {
	VerifyUser *VerifyUser
}

type CreateTokenResponse struct {
	VerifiedUser *VerifiedUser
}
