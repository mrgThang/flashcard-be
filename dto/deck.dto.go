package dto

type CreateDeckRequest struct {
	Name        string
	Description string
	UserID      int32
}

type UpdateDeckRequest struct {
	ID          int32
	Name        string
	Description string
}

type GetDeckRequest struct {
	ID int32
}
