package dto

type CreateDeckRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int32
}

type UpdateDeckRequest struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetDecksRequest struct {
	ID       int32
	Name     string
	UserID   int32
	Page     int
	PageSize int
}

type GetDecksResponse struct {
	Pagination
	Decks []DeckItem `json:"decks"`
}

type DeckItem struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TotalCards  int32  `json:"totalCards"`
	CardsLeft   int32  `json:"cardsLeft"`
}
