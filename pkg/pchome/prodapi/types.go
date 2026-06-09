package prodapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type BoolInt int

const (
	Bool0 BoolInt = 0
	Bool1 BoolInt = 1
)

func (b *BoolInt) UnmarshalJSON(data []byte) error {
	value, err := parseOptionalInt(data, "BoolInt")
	if err != nil {
		return err
	}
	if value == nil {
		*b = Bool0
		return nil
	}
	*b = BoolInt(*value)
	return nil
}

type Price struct {
	M     int    `json:"M"`
	P     int    `json:"P"`
	Prime string `json:"Prime"`
	Low   *int   `json:"Low"`
}

func (p *Price) UnmarshalJSON(data []byte) error {
	var raw struct {
		M     int             `json:"M"`
		P     int             `json:"P"`
		Prime string          `json:"Prime"`
		Low   json.RawMessage `json:"Low"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	low, err := parseOptionalInt(raw.Low, "Price.Low")
	if err != nil {
		return err
	}

	p.M = raw.M
	p.P = raw.P
	p.Prime = raw.Prime
	p.Low = low
	return nil
}

func parseOptionalInt(data json.RawMessage, field string) (*int, error) {
	text := strings.TrimSpace(string(data))
	if text == "" || text == "null" {
		return nil, nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("%s must be an integer string: %w", field, err)
		}
		return &n, nil
	}

	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		return &n, nil
	}

	return nil, fmt.Errorf("%s must be an integer, integer string, or null", field)
}

type Picture struct {
	B string `json:"B"`
	S string `json:"S"`
	W string `json:"W"`
}

type Attribute struct {
	Id      string   `json:"Id"`
	Options []string `json:"Options"`
}

type Product struct {
	Id          string      `json:"Id"`
	Name        string      `json:"Name"`
	Nick        string      `json:"Nick"`
	SalesName   string      `json:"SalesName"`
	Tagline     string      `json:"Tagline"`
	Qty         *int        `json:"Qty"`
	ShipDay     int         `json:"ShipDay"`
	ShipType    string      `json:"ShipType"`
	Price       *Price      `json:"Price"`
	Pic         *Picture    `json:"Pic"`
	PicExtra    []string    `json:"PicExtra"`
	BrandList   []string    `json:"BrandList"`
	CategoryIds []string    `json:"CategoryIds"`
	Attributes  []Attribute `json:"Attributes"`

	IsArrival24h    *BoolInt `json:"isArrival24h"`
	IsOnSale        *BoolInt `json:"isOnSale"`
	IsOrderDiscount *BoolInt `json:"isOrderDiscount"`
	IsPrimeOnly     *BoolInt `json:"isPrimeOnly"`

	RatingValue *float64 `json:"RatingValue"`
	ReviewCount *int     `json:"ReviewCount"`
}

func (p *Product) UnmarshalJSON(data []byte) error {
	type productAlias Product
	var raw struct {
		*productAlias
		Qty         json.RawMessage `json:"Qty"`
		ShipDay     json.RawMessage `json:"ShipDay"`
		ReviewCount json.RawMessage `json:"ReviewCount"`
	}
	raw.productAlias = (*productAlias)(p)

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	qty, err := parseOptionalInt(raw.Qty, "Product.Qty")
	if err != nil {
		return err
	}
	shipDay, err := parseIntDefault(raw.ShipDay, "Product.ShipDay")
	if err != nil {
		return err
	}
	reviewCount, err := parseOptionalInt(raw.ReviewCount, "Product.ReviewCount")
	if err != nil {
		return err
	}

	p.Qty = qty
	p.ShipDay = shipDay
	p.ReviewCount = reviewCount
	return nil
}

type ButtonItem struct {
	Seq   int    `json:"Seq"`
	Id    string `json:"Id"`
	Store string `json:"Store"`
	Group string `json:"Group"`

	Price *Price `json:"Price"`
	Qty   *int   `json:"Qty"`

	ButtonType string `json:"ButtonType"`
	SaleStatus int    `json:"SaleStatus"`
	SpecialQty int    `json:"SpecialQty"`
	Device     any    `json:"Device"`

	IsOrderDiscount *BoolInt `json:"isOrderDiscount"`
	IsPrimeOnly     *BoolInt `json:"isPrimeOnly"`
}

func (b *ButtonItem) UnmarshalJSON(data []byte) error {
	type buttonAlias ButtonItem
	var raw struct {
		*buttonAlias
		Seq        json.RawMessage `json:"Seq"`
		Qty        json.RawMessage `json:"Qty"`
		SaleStatus json.RawMessage `json:"SaleStatus"`
		SpecialQty json.RawMessage `json:"SpecialQty"`
	}
	raw.buttonAlias = (*buttonAlias)(b)

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	seq, err := parseIntDefault(raw.Seq, "ButtonItem.Seq")
	if err != nil {
		return err
	}
	qty, err := parseOptionalInt(raw.Qty, "ButtonItem.Qty")
	if err != nil {
		return err
	}
	saleStatus, err := parseIntDefault(raw.SaleStatus, "ButtonItem.SaleStatus")
	if err != nil {
		return err
	}
	specialQty, err := parseIntDefault(raw.SpecialQty, "ButtonItem.SpecialQty")
	if err != nil {
		return err
	}

	b.Seq = seq
	b.Qty = qty
	b.SaleStatus = saleStatus
	b.SpecialQty = specialQty
	return nil
}

func parseIntDefault(data json.RawMessage, field string) (int, error) {
	value, err := parseOptionalInt(data, field)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return 0, nil
	}
	return *value, nil
}
