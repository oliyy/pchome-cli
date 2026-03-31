package cmd

import (
	"bytes"
	"testing"

	"github.com/oliy/pchome-cli/pkg/i18n"
)

func TestRenderSearchResult_Text(t *testing.T) {
	var buf bytes.Buffer
	if err := renderSearchResult(&buf, i18n.New(i18n.LangZHTW), sampleSearchResult(), "#,price,rating,reviews,24h,stock,brand,name,id", false, 60); err != nil {
		t.Fatalf("renderSearchResult returned error: %v", err)
	}

	assertGolden(t, "render/search_zh_tw.golden", buf.String())
}

func TestRenderProductDetail_Text(t *testing.T) {
	var buf bytes.Buffer
	renderProductDetail(&buf, i18n.New(i18n.LangZHTW), sampleProductDetail())

	assertGolden(t, "render/product_detail_zh_tw.golden", buf.String())
}
