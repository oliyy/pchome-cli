package hermes

// Hermes recommendations API (observed shapes).
//
// Endpoint:
//   POST https://apih.pcloud.tw/hermes/api/goods/rank
//
// Notes:
// - `token` is required; missing/invalid tokens return an error JSON object.
// - Names can include HTML entities; use html.UnescapeString when displaying.

type GoodsRankRequest struct {
	Token   string `json:"token"`
	RecType string `json:"rec_type"`
	RecPos  string `json:"rec_pos"`
	TopK    int    `json:"topk"`
	Device  string `json:"device"`
	UID     any    `json:"uid"`

	VenSession string `json:"ven_session,omitempty"`
	VenGuid    string `json:"ven_guid,omitempty"`

	GID       any `json:"gid,omitempty"`
	CategCode any `json:"categ_code,omitempty"`

	BCategInfo []string               `json:"b_categ_info,omitempty"`
	WCategInfo []GoodsRankCategory    `json:"w_categ_info,omitempty"`
	RefInfo    []GoodsRankRefInfoItem `json:"ref_info,omitempty"`
}

type GoodsRankCategory struct {
	Code string `json:"code"`
}

type GoodsRankRefInfoItem struct {
	CategCode string `json:"categ_code"`
}

type GoodsRankRefItem struct {
	ItemName string  `json:"item_name"`
	AlgName  string  `json:"alg_name"`
	RefModel string  `json:"ref_model"`
	Score    float64 `json:"score"`
}

type GoodsRankItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	SalePrice   float64 `json:"sale_price"`
	GoodsImgURL string  `json:"goods_img_url"`

	MsgType string `json:"msg_type"`
	Msg     string `json:"msg"`

	RefItemList []GoodsRankRefItem `json:"ref_item_list"`
}

type GoodsRankResponse struct {
	RecomdID   string          `json:"recomd_id"`
	Took       float64         `json:"took"`
	TimedOut   bool            `json:"timed_out"`
	TotalHits  float64         `json:"total_hits"`
	RecomdList []GoodsRankItem `json:"recomd_list"`

	// Error shape
	Input           any    `json:"input"`
	Error           string `json:"error"`
	BackendResponse string `json:"back-end response"`
}
