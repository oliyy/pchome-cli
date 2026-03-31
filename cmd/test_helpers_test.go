package cmd

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/oliy/pchome-cli/pkg/catalog"
)

var updateGolden = flag.Bool("update", false, "update golden test files")

func assertGolden(t *testing.T, relativePath, got string) {
	t.Helper()

	path := filepath.Join("testdata", relativePath)
	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create golden dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden file: %v", err)
		}
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden file %s: %v", path, err)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", relativePath, got, string(want))
	}
}

func runRootHelp(t *testing.T, configData string) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)

	if configData != "" {
		configDir := filepath.Join(home, ".pchome")
		if err := os.MkdirAll(configDir, 0o700); err != nil {
			t.Fatalf("create config dir: %v", err)
		}
		configPath := filepath.Join(configDir, "config.toml")
		if err := os.WriteFile(configPath, []byte(configData), 0o600); err != nil {
			t.Fatalf("write config: %v", err)
		}
	}

	root, err := newRootCmd()
	if err != nil {
		t.Fatalf("newRootCmd returned error: %v", err)
	}

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	return buf.String()
}

func sampleSearchResult() *catalog.SearchResult {
	rating := 4.9
	reviews := 33
	stock := 20
	arrival := true
	inStock := true

	return &catalog.SearchResult{
		Query:        "掃地機器人",
		Sort:         "relevance",
		Page:         1,
		PagesScanned: 1,
		TotalRows:    100,
		Returned:     1,
		Items: []catalog.ProductSummary{
			{
				ID:    "DMBL53-A900JDNJS",
				Name:  "小米 Xiaomi 掃拖機器人H40",
				Brand: "Xiaomi",
				Price: catalog.Money{Current: 6999},
				Availability: catalog.Availability{
					InStock:    &inStock,
					StockQty:   &stock,
					Arrival24h: &arrival,
				},
				Rating: catalog.Rating{
					Score:   &rating,
					Reviews: &reviews,
				},
				URL: "https://24h.pchome.com.tw/prod/DMBL53-A900JDNJS",
			},
		},
	}
}

func sampleProductDetail() *catalog.ProductDetail {
	rating := 5.0
	reviews := 4
	stock := 9
	arrival := true

	return &catalog.ProductDetail{
		Product: catalog.ProductSummary{
			ID:    "DRAA5K-A900JOK9O",
			Name:  "Micron Crucial X10",
			Brand: "Micron",
			Price: catalog.Money{Current: 4340, Original: 5999, DiscountPercent: 28},
			Availability: catalog.Availability{
				StockQty:   &stock,
				Arrival24h: &arrival,
				ShipType:   "Consign",
			},
			Rating: catalog.Rating{
				Score:   &rating,
				Reviews: &reviews,
			},
			URL:      "https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O",
			ImageURL: "https://example.com/image.jpg",
		},
		SalesName:    "Crucial X10 1TB",
		BrandAliases: []string{"Micron", "美光"},
		Images:       []string{"https://example.com/image.jpg"},
	}
}
