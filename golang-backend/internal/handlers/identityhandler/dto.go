package identityhandler

// TODO: this struct is duplicated on user handler
type userPermissions struct {
	IsSuperUser bool `json:"isSuperUser"`
}

type userResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	CreatedAt   int64           `json:"createdAt"`
	UpdatedAt   *int64          `json:"updatedAt"`
	DeletedAt   *int64          `json:"deletedAt"`
	Permissions userPermissions `json:"permissions"`
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

type signInResponse struct {
	AccessToken  tokenResponse `json:"accessToken"`
	RefreshToken tokenResponse `json:"refreshToken"`
	User         userResponse  `json:"user"`
}

type renewAccessTokenResponse struct {
	AccessToken tokenResponse `json:"accessToken"`
	User        userResponse  `json:"user"`
}
