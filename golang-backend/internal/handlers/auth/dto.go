package handlers

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"passsword"`
}

type SuccessSignInResponse struct {
	Token string `json:"token"`
}
