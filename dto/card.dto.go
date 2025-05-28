package dto

type CreateCardRequest struct {
	Front  string `json:"front"`
	Back   string `json:"back"`
	DeckID int32  `json:"deckId"`
	UserID int32  `json:"userId"`
}

type UpdateCardRequest struct {
	ID    int32  `json:"id"`
	Front string `json:"front"`
	Back  string `json:"back"`
}

type GetCardsRequest struct {
	ID       int32  `json:"id"`
	DeckID   int32  `json:"deckId"`
	UserID   int32  `json:"userId"`
	Front    string `json:"front"`
	Back     string `json:"back"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
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
	UserID int32  `json:"userId"`
}
