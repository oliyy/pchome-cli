package prodapi

type BoolInt int

const (
	Bool0 BoolInt = 0
	Bool1 BoolInt = 1
)

type Price struct {
	M     int    `json:"M"`
	P     int    `json:"P"`
	Prime string `json:"Prime"`
	Low   *int   `json:"Low"`
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
