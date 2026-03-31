package search

// Types are intentionally partial: only fields used by the CLI are strongly typed.
// The underlying endpoints are undocumented and may change.

type ResultsResponse struct {
	QTime     int             `json:"QTime"`
	TotalRows int             `json:"TotalRows"`
	TotalPage int             `json:"TotalPage"`
	Q         string          `json:"q"`
	Prods     []ResultProduct `json:"Prods"`
}

type ResultProduct struct {
	Id          string   `json:"Id"`
	Name        string   `json:"Name"`
	Describe    string   `json:"Describe"`
	Brand       string   `json:"Brand"`
	Price       int      `json:"Price"`
	OriginPrice int      `json:"OriginPrice"`
	PCateID     []string `json:"PCateId"`
	PicS        string   `json:"PicS"`
	PicB        string   `json:"PicB"`
	PublishDate string   `json:"PublishDate"`

	RatingValue *float64 `json:"ratingValue"`
	ReviewCount *int     `json:"reviewCount"`
}

type CategoryNode struct {
	Id    string         `json:"Id"`
	Name  *string        `json:"Name,omitempty"`
	Qty   int            `json:"Qty"`
	Nodes []CategoryNode `json:"Nodes,omitempty"`
}

type PCategoryNode struct {
	Id    string          `json:"Id"`
	Name  string          `json:"Name"`
	Qty   int             `json:"Qty"`
	Nodes []PCategoryNode `json:"Nodes,omitempty"`
}

type BrandFacet struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
	Qty  int    `json:"Qty"`
}

type GroupAttributeGroup struct {
	Id              string                 `json:"Id"`
	Name            string                 `json:"Name"`
	Qty             int                    `json:"Qty"`
	CategoryOptions []GroupAttributeBucket `json:"CategoryOptions"`
	OtherOptions    []GroupAttributeBucket `json:"OtherOptions"`
}

type GroupAttributeBucket struct {
	CategoryId   string                 `json:"CategoryId"`
	CategoryName string                 `json:"CategoryName"`
	Options      []GroupAttributeOption `json:"Options"`
}

type GroupAttributeOption struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
	Qty  int    `json:"Qty"`
}

type SuggestWord struct {
	Name string `json:"Name"`
	Rurl string `json:"Rurl"`
	Src  string `json:"Src"`
	BU   string `json:"BU"`
}
