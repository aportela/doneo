package handlers

type UserRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	IsSuperUser bool   `json:"isSuperUser"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   *int64 `json:"updatedAt,omitempty"`
	IsSuperUser bool   `json:"isSuperUser"`
	AvatarURL   string `json:"avatar"`
}

type GetUserResponse struct {
	User UserResponse `json:"user"`
}

type SearchUserResponse struct {
	Users []UserResponse `json:"users"`
}
