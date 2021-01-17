package models

type Offer struct {
	SellerID int64  `json:"seller_id"`
	OfferID  int    `json:"offer_id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}
