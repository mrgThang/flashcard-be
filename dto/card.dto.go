package dto

type CreateCardRequest struct {
	Front  string
	Back   string
	DeckID int32
	UserID int32
}

type UpdateCardRequest struct {
	ID    int32
	Front string
	Back  string
}

type GetCardsRequest struct {
	ID     int32
	DeckID int32
	UserID int32
	Front  string
	Back   string
	OffSet int
	Limit  int
}
