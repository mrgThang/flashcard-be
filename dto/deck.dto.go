package dto

type CreateDeckRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int32  `json:"userId"`
}

type UpdateDeckRequest struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetDecksRequest struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	UserID   int32  `json:"userId"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type GetDecksResponse struct {
	Pagination
	Decks []DeckItem `json:"decks"`
}

type DeckItem struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	UserID int32  `json:"userId"`
}
