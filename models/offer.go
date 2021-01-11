package models

type Offer struct {
	SellerID int64  `json:"seller_id"`
	OfferID  int64  `json:"offer_id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int64  `json:"quantity"`
}
