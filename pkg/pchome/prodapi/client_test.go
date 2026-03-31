package prodapi

import "testing"

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
