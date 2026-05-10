package authhandler

type userResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   *int64 `json:"updatedAt"`
	DeletedAt   *int64 `json:"deletedAt"`
	IsSuperUser bool   `json:"isSuperUser"`
	AvatarURL   string `json:"avatarUrl"`
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

type SignInResponse struct {
	AccessToken  TokenResponse `json:"accessToken"`
	RefreshToken TokenResponse `json:"refreshToken"`
	User         userResponse  `json:"user"`
}

type RenewAccessTokenResponse struct {
	AccessToken TokenResponse `json:"accessToken"`
	User        userResponse  `json:"user"`
}
