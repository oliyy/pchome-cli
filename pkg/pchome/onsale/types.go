package onsale

// OnSale CMS payload (observed shapes).
//
// Endpoint:
//   GET https://ecapi-cdn.pchome.com.tw/fsapi/cms/onsale

type Response struct {
	Data []Slot `json:"data"`
}

type Slot struct {
	Slot      string `json:"slot"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Status    string `json:"status"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`

	Products []Product `json:"products"`
}

type Product struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slogan string `json:"slogan"`
	URL    string `json:"url"`

	Sorted          *int   `json:"sorted"`
	ImageSourceType string `json:"imageSourceType"`
	Image           string `json:"image"`

	Price Price `json:"price"`
}

type Price struct {
	Origin   int     `json:"origin"`
	OnSale   int     `json:"onsale"`
	Unit     *string `json:"unit"`
	Discount int     `json:"discount"`
}
