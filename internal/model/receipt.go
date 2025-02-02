package model

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"` // Keep as string
	PurchaseTime string `json:"purchaseTime"` // Keep as string
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"` // Keep as string
}

type PointsResponse struct {
	Points int `json:"points"`
}

type IdResponse struct {
	ID string `json:"id"`
}
