package dto

type GetUserRequest struct {
	ID int32 `json:"id"`
}

type GetUserResponse struct {
	User UserItem `json:"user"`
}

type UserItem struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID int32 `json:"id"`
}

type UpdateUserRequest struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type UpdateUserResponse struct {
	ID int32 `json:"id"`
}
