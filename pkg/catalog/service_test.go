package catalog

import (
	"testing"

	"github.com/oliy/pchome-cli/pkg/pchome/hermes"
	"github.com/oliy/pchome-cli/pkg/pchome/prodapi"
	"github.com/oliy/pchome-cli/pkg/pchome/search"
)

func TestSummaryFromRecordPrefersEnrichedProductData(t *testing.T) {
	rating := 4.9
	reviews := 12
	low := 3999
	m := prodapi.Bool1
	stock := 7
	inStock := true

	rec := &productRecord{
		id: "DRAA5K-A900JOK9O",
		searchProduct: &search.ResultProduct{
			Id:          "DRAA5K-A900JOK9O",
			Name:        "Search Name",
			Describe:    "Search description",
			Brand:       "Search Brand",
			Price:       4999,
			OriginPrice: 5999,
			PCateID:     []string{"A"},
			PicB:        "/items/SEARCH/000001.jpg",
		},
		product: &prodapi.Product{
			Id:          "DRAA5K-A900JOK9O-000",
			Name:        "Full Product Name",
			Tagline:     "Detailed tagline",
			BrandList:   []string{"Micron", "美光"},
			CategoryIds: []string{"A", "B"},
			Price: &prodapi.Price{
				M:   5999,
				P:   4540,
				Low: &low,
			},
			IsArrival24h: &m,
			ShipType:     "Consign",
			RatingValue:  &rating,
			ReviewCount:  &reviews,
		},
		stockQty: &stock,
		inStock:  &inStock,
	}

	summary := summaryFromRecord(rec)

	if summary.Name != "Full Product Name" {
		t.Fatalf("expected enriched product name, got %q", summary.Name)
	}
	if summary.Brand != "Micron" {
		t.Fatalf("expected brand Micron, got %q", summary.Brand)
	}
	if summary.Price.Current != 3999 {
		t.Fatalf("expected current price 3999, got %d", summary.Price.Current)
	}
	if summary.Price.Original != 5999 {
		t.Fatalf("expected original price 5999, got %d", summary.Price.Original)
	}
	if summary.Price.DiscountPercent == 0 {
		t.Fatalf("expected discount percent to be derived")
	}
	if summary.Availability.Arrival24h == nil || !*summary.Availability.Arrival24h {
		t.Fatalf("expected arrival_24h to be true")
	}
	if summary.Availability.StockQty == nil || *summary.Availability.StockQty != 7 {
		t.Fatalf("expected stock qty 7")
	}
	if len(summary.CategoryIDs) != 2 {
		t.Fatalf("expected merged category ids, got %v", summary.CategoryIDs)
	}
}

func TestHermesTokenFallsBackToBundledDefault(t *testing.T) {
	t.Setenv("PCHOME_HERMES_TOKEN", "")

	if got := hermesToken(""); got != hermes.DefaultToken {
		t.Fatalf("expected bundled fallback token, got %q", got)
	}
}

func TestHermesTokenPrefersConfigOverEnv(t *testing.T) {
	t.Setenv("PCHOME_HERMES_TOKEN", "from-env")

	if got := hermesToken("from-config"); got != "from-config" {
		t.Fatalf("expected config token precedence, got %q", got)
	}
	if got := hermesToken(""); got != "from-env" {
		t.Fatalf("expected env token fallback, got %q", got)
	}
}
