package userauthnvalue

type VerifyUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifiedUser struct {
	AccessToken string `json:"access_token"`
}

type CreateTokenRequest struct {
	VerifyUser *VerifyUser `json:"verify_user"`
}

type CreateTokenResponse struct {
	VerifiedUser *VerifiedUser `json:"verified_user"`
}
