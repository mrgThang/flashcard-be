package dto

type CreateCardRequest struct {
	Front  string `json:"front"`
	Back   string `json:"back"`
	DeckID int32  `json:"deckId"`
	UserID int32
}

type UpdateCardRequest struct {
	ID    int32  `json:"id"`
	Front string `json:"front"`
	Back  string `json:"back"`
}

type GetCardsRequest struct {
	ID       int32
	DeckID   int32
	UserID   int32
	Front    string
	Back     string
	Page     int
	PageSize int
}

type GetCardsResponse struct {
	Pagination
	Cards []CardItem `json:"cards"`
}

type CardItem struct {
	ID     int32  `json:"id"`
	Front  string `json:"front"`
	Back   string `json:"back"`
	DeckID int32  `json:"deckId"`
}
