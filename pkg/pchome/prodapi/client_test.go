package prodapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestExtractFirstJSONPArgument(t *testing.T) {
	in := `try{jsonp_cb({"a":1});}catch(e){if(window.console){console.log(e);}}`
	got, ok := extractFirstJSONPArgument(in)
	if !ok {
		t.Fatalf("expected ok")
	}
	if got != `{"a":1}` {
		t.Fatalf("unexpected payload: %q", got)
	}
}

func TestExtractFirstJSONPArgument_IgnoresParensInsideStrings(t *testing.T) {
	in := `try{cb({"a":"(x)","b":"y)z"});}catch(e){}`
	got, ok := extractFirstJSONPArgument(in)
	if !ok {
		t.Fatalf("expected ok")
	}
	if got != `{"a":"(x)","b":"y)z"}` {
		t.Fatalf("unexpected payload: %q", got)
	}
}

func TestParseJSONOrJSONP(t *testing.T) {
	plain := []byte(`{"x":1}`)
	got, err := parseJSONOrJSONP(plain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != `{"x":1}` {
		t.Fatalf("unexpected plain: %q", string(got))
	}

	jsonp := []byte(`try{cb([{"x":1}]);}catch(e){}`)
	got, err = parseJSONOrJSONP(jsonp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != `[{"x":1}]` {
		t.Fatalf("unexpected jsonp: %q", string(got))
	}
}

func TestDecodeProductsPayload_PriceLowVariants(t *testing.T) {
	payload := []byte(`{
		"DMABP5-A900I86D3-000": {"Id":"DMABP5-A900I86D3-000","Price":{"M":1980,"P":968,"Low":null,"Prime":""}},
		"DMABLY-A900INHLM-000": {"Id":"DMABLY-A900INHLM-000","Price":{"M":3990,"P":1388,"Low":"1288","Prime":""}},
		"DMAB46-A900HFB5Y-000": {"Id":"DMAB46-A900HFB5Y-000","Qty":"5","ShipDay":"0","ReviewCount":"508","Price":{"M":1288,"P":699,"Low":594,"Prime":""}}
	}`)

	got, err := decodeProductsPayload(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DMABP5-A900I86D3-000"].Price.Low != nil {
		t.Fatalf("expected null Low to stay nil")
	}
	if low := got["DMABLY-A900INHLM-000"].Price.Low; low == nil || *low != 1288 {
		t.Fatalf("expected string Low 1288, got %#v", low)
	}
	if low := got["DMAB46-A900HFB5Y-000"].Price.Low; low == nil || *low != 594 {
		t.Fatalf("expected numeric Low 594, got %#v", low)
	}
	if qty := got["DMAB46-A900HFB5Y-000"].Qty; qty == nil || *qty != 5 {
		t.Fatalf("expected string Qty 5, got %#v", qty)
	}
	if got["DMAB46-A900HFB5Y-000"].ShipDay != 0 {
		t.Fatalf("expected string ShipDay 0, got %d", got["DMAB46-A900HFB5Y-000"].ShipDay)
	}
	if reviews := got["DMAB46-A900HFB5Y-000"].ReviewCount; reviews == nil || *reviews != 508 {
		t.Fatalf("expected string ReviewCount 508, got %#v", reviews)
	}
}

func TestDecodeProductsPayload_ArrayShapes(t *testing.T) {
	empty, err := decodeProductsPayload([]byte(`[]`))
	if err != nil {
		t.Fatalf("unexpected empty array error: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected empty map, got %d products", len(empty))
	}

	got, err := decodeProductsPayload([]byte(`[{"Id":"DMABLY-A900INHLM-000","Price":{"Low":"1288"}}]`))
	if err != nil {
		t.Fatalf("unexpected product array error: %v", err)
	}
	if _, ok := got["DMABLY-A900INHLM"]; !ok {
		t.Fatalf("expected normalized id key, got %#v", got)
	}
}

func TestDecodeProductsPayload_PreservesDecodeError(t *testing.T) {
	_, err := decodeProductsPayload([]byte(`{"D-000":{"Id":"D-000","Price":{"Low":"sale"}}}`))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "Price.Low") {
		t.Fatalf("expected field-level error, got %v", err)
	}
}

func TestButtonItemIntegerVariants(t *testing.T) {
	var got ButtonItem
	err := json.Unmarshal([]byte(`{
		"Seq":"1",
		"Id":"DMABN6-A900JUFY7-000",
		"Qty":"5",
		"SaleStatus":"1",
		"SpecialQty":"2",
		"ButtonType":"ForSale",
		"isPrimeOnly":"0",
		"isOrderDiscount":"1"
	}`), &got)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Seq != 1 {
		t.Fatalf("expected Seq 1, got %d", got.Seq)
	}
	if got.Qty == nil || *got.Qty != 5 {
		t.Fatalf("expected Qty 5, got %#v", got.Qty)
	}
	if got.SaleStatus != 1 {
		t.Fatalf("expected SaleStatus 1, got %d", got.SaleStatus)
	}
	if got.SpecialQty != 2 {
		t.Fatalf("expected SpecialQty 2, got %d", got.SpecialQty)
	}
	if got.IsPrimeOnly == nil || *got.IsPrimeOnly != Bool0 {
		t.Fatalf("expected isPrimeOnly 0, got %#v", got.IsPrimeOnly)
	}
	if got.IsOrderDiscount == nil || *got.IsOrderDiscount != Bool1 {
		t.Fatalf("expected isOrderDiscount 1, got %#v", got.IsOrderDiscount)
	}
}
