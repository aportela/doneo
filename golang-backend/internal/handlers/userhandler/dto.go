package userhandler

type addRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsSuperUser bool   `json:"isSuperUser"`
}

type updateRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    *string `json:"password,omitempty"`
	IsSuperUser bool    `json:"isSuperUser"`
}

type userResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   *int64 `json:"updatedAt"`
	DeletedAt   *int64 `json:"deletedAt"`
	IsSuperUser bool   `json:"isSuperUser"`
	AvatarURL   string `json:"avatar"`
}

type addResponse struct {
	User userResponse `json:"user"`
}

type updateResponse struct {
	User userResponse `json:"user"`
}

type getResponse struct {
	User userResponse `json:"user"`
}

type searchResponse struct {
	Users []userResponse `json:"users"`
}
