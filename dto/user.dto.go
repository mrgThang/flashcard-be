package dto

type GetUserRequest struct {
	ID uint
}

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserRequest struct {
	ID       uint
	Name     string
	Email    string
	Password string
}
