package exturi_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/exturi"
)

// ---------- URIParse ----------

func TestURIParse_Full(t *testing.T) {
	fn := exturi.URIParse()
	got, err := fn([]any{"https://user:pass@example.com:8080/path/to?a=1&b=2#frag"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	check := func(key, want string) {
		t.Helper()
		if obj[key] != want {
			t.Errorf("%s: expected %q, got %q", key, want, obj[key])
		}
	}
	check("scheme", "https")
	check("user", "user")
	check("password", "pass")
	check("host", "example.com")
	check("port", "8080")
	check("path", "/path/to")
	check("fragment", "frag")
	q := obj["query"].(map[string]any)
	if q["a"] != "1" {
		t.Errorf("query.a: expected '1', got %v", q["a"])
	}
}

func TestURIParse_Simple(t *testing.T) {
	fn := exturi.URIParse()
	got, err := fn([]any{"https://example.com"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["host"] != "example.com" {
		t.Errorf("host: expected 'example.com', got %v", obj["host"])
	}
	if obj["path"] != "" {
		t.Errorf("path: expected '', got %v", obj["path"])
	}
}

func TestURIParse_NoArgs(t *testing.T) {
	fn := exturi.URIParse()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for no args")
	}
}

func TestURIParse_WrongType(t *testing.T) {
	fn := exturi.URIParse()
	_, err := fn([]any{123}, nil)
	if err == nil {
		t.Error("expected error for non-string")
	}
}

func TestURIParse_NoUserInfo(t *testing.T) {
	fn := exturi.URIParse()
	got, err := fn([]any{"http://example.com/"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["user"] != "" {
		t.Errorf("user: expected '', got %v", obj["user"])
	}
}

// ---------- URIBuild ----------

func TestURIBuild_Basic(t *testing.T) {
	fn := exturi.URIBuild()
	parts := map[string]any{
		"scheme": "https",
		"host":   "example.com",
		"path":   "/api/v1",
	}
	got, err := fn([]any{parts}, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := got.(string)
	if s != "https://example.com/api/v1" {
		t.Errorf("expected 'https://example.com/api/v1', got %q", s)
	}
}

func TestURIBuild_WithPort(t *testing.T) {
	fn := exturi.URIBuild()
	parts := map[string]any{
		"scheme": "http",
		"host":   "localhost",
		"port":   "9090",
		"path":   "/",
	}
	got, _ := fn([]any{parts}, nil)
	s := got.(string)
	if s != "http://localhost:9090/" {
		t.Errorf("expected 'http://localhost:9090/', got %q", s)
	}
}

func TestURIBuild_WithQuery(t *testing.T) {
	fn := exturi.URIBuild()
	parts := map[string]any{
		"scheme": "https",
		"host":   "example.com",
		"path":   "/search",
		"query":  map[string]any{"q": "hello"},
	}
	got, err := fn([]any{parts}, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := got.(string)
	if s != "https://example.com/search?q=hello" {
		t.Errorf("unexpected: %q", s)
	}
}

func TestURIBuild_WithUserPassword(t *testing.T) {
	fn := exturi.URIBuild()
	parts := map[string]any{
		"scheme":   "ftp",
		"user":     "alice",
		"password": "secret",
		"host":     "ftp.example.com",
	}
	got, _ := fn([]any{parts}, nil)
	s := got.(string)
	if s != "ftp://alice:secret@ftp.example.com" {
		t.Errorf("unexpected: %q", s)
	}
}

func TestURIBuild_WithUserOnly(t *testing.T) {
	fn := exturi.URIBuild()
	parts := map[string]any{
		"scheme": "ftp",
		"user":   "alice",
		"host":   "ftp.example.com",
	}
	got, _ := fn([]any{parts}, nil)
	s := got.(string)
	if s != "ftp://alice@ftp.example.com" {
		t.Errorf("unexpected: %q", s)
	}
}

func TestURIBuild_NoArgs(t *testing.T) {
	fn := exturi.URIBuild()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for no args")
	}
}

func TestURIBuild_WrongType(t *testing.T) {
	fn := exturi.URIBuild()
	_, err := fn([]any{"not-an-object"}, nil)
	if err == nil {
		t.Error("expected error for non-object")
	}
}

// ---------- URIJoin ----------

func TestURIJoin_Relative(t *testing.T) {
	fn := exturi.URIJoin()
	got, err := fn([]any{"https://example.com/a/b/", "../c"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://example.com/a/c" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestURIJoin_Absolute(t *testing.T) {
	fn := exturi.URIJoin()
	got, _ := fn([]any{"https://example.com/a/b", "https://other.com/x"}, nil)
	if got != "https://other.com/x" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestURIJoin_NoArgs(t *testing.T) {
	fn := exturi.URIJoin()
	_, err := fn([]any{"base"}, nil)
	if err == nil {
		t.Error("expected error for missing ref")
	}
}

func TestURIJoin_WrongType(t *testing.T) {
	fn := exturi.URIJoin()
	_, err := fn([]any{123, "ref"}, nil)
	if err == nil {
		t.Error("expected error for non-string base")
	}
}

// ---------- QueryParse ----------

func TestQueryParse_Basic(t *testing.T) {
	fn := exturi.QueryParse()
	got, err := fn([]any{"a=1&b=hello"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["a"] != "1" {
		t.Errorf("a: expected '1', got %v", obj["a"])
	}
	if obj["b"] != "hello" {
		t.Errorf("b: expected 'hello', got %v", obj["b"])
	}
}

func TestQueryParse_WithLeadingQuestion(t *testing.T) {
	fn := exturi.QueryParse()
	got, err := fn([]any{"?x=10"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["x"] != "10" {
		t.Errorf("x: expected '10', got %v", obj["x"])
	}
}

func TestQueryParse_MultiValue(t *testing.T) {
	fn := exturi.QueryParse()
	got, err := fn([]any{"tag=a&tag=b&tag=c"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	arr, ok := obj["tag"].([]any)
	if !ok {
		t.Fatalf("tag: expected []any, got %T", obj["tag"])
	}
	if len(arr) != 3 {
		t.Errorf("tag: expected 3 values, got %d", len(arr))
	}
}

func TestQueryParse_Empty(t *testing.T) {
	fn := exturi.QueryParse()
	got, err := fn([]any{""}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if len(obj) != 0 {
		t.Errorf("expected empty object, got %v", obj)
	}
}

func TestQueryParse_NoArgs(t *testing.T) {
	fn := exturi.QueryParse()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for no args")
	}
}

func TestQueryParse_WrongType(t *testing.T) {
	fn := exturi.QueryParse()
	_, err := fn([]any{42}, nil)
	if err == nil {
		t.Error("expected error for non-string")
	}
}

// ---------- QueryBuild ----------

func TestQueryBuild_Basic(t *testing.T) {
	fn := exturi.QueryBuild()
	got, err := fn([]any{map[string]any{"a": "1", "b": "hello"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := got.(string)
	// Sorted keys: a=1&b=hello
	if s != "a=1&b=hello" {
		t.Errorf("unexpected: %q", s)
	}
}

func TestQueryBuild_ArrayValue(t *testing.T) {
	fn := exturi.QueryBuild()
	got, err := fn([]any{map[string]any{"tag": []any{"a", "b"}}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := got.(string)
	if s != "tag=a&tag=b" {
		t.Errorf("unexpected: %q", s)
	}
}

func TestQueryBuild_NoArgs(t *testing.T) {
	fn := exturi.QueryBuild()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for no args")
	}
}

func TestQueryBuild_WrongType(t *testing.T) {
	fn := exturi.QueryBuild()
	_, err := fn([]any{"not-a-map"}, nil)
	if err == nil {
		t.Error("expected error for non-object")
	}
}

// ---------- URIGetPath ----------

func TestURIGetPath(t *testing.T) {
	fn := exturi.URIGetPath()
	got, err := fn([]any{"https://example.com/foo/bar?x=1"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "/foo/bar" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestURIGetPath_NoArgs(t *testing.T) {
	fn := exturi.URIGetPath()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestURIGetPath_WrongType(t *testing.T) {
	fn := exturi.URIGetPath()
	_, err := fn([]any{99}, nil)
	if err == nil {
		t.Error("expected error for non-string")
	}
}

// ---------- URIGetQuery ----------

func TestURIGetQuery(t *testing.T) {
	fn := exturi.URIGetQuery()
	got, err := fn([]any{"https://example.com/?a=1&b=2"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["a"] != "1" || obj["b"] != "2" {
		t.Errorf("unexpected: %v", obj)
	}
}

func TestURIGetQuery_NoArgs(t *testing.T) {
	fn := exturi.URIGetQuery()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestURIGetQuery_WrongType(t *testing.T) {
	fn := exturi.URIGetQuery()
	_, err := fn([]any{false}, nil)
	if err == nil {
		t.Error("expected error for non-string")
	}
}

// ---------- URISetQuery ----------

func TestURISetQuery(t *testing.T) {
	fn := exturi.URISetQuery()
	got, err := fn([]any{"https://example.com/search?old=1", map[string]any{"q": "new"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := got.(string)
	if s != "https://example.com/search?q=new" {
		t.Errorf("unexpected: %q", s)
	}
}

func TestURISetQuery_EmptyParams(t *testing.T) {
	fn := exturi.URISetQuery()
	got, err := fn([]any{"https://example.com/path?a=1", map[string]any{}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := got.(string)
	if s != "https://example.com/path" {
		t.Errorf("unexpected: %q", s)
	}
}

func TestURISetQuery_NoArgs(t *testing.T) {
	fn := exturi.URISetQuery()
	_, err := fn([]any{"url"}, nil)
	if err == nil {
		t.Error("expected error for missing params")
	}
}

func TestURISetQuery_WrongURLType(t *testing.T) {
	fn := exturi.URISetQuery()
	_, err := fn([]any{42, map[string]any{}}, nil)
	if err == nil {
		t.Error("expected error for non-string url")
	}
}

func TestURISetQuery_WrongParamsType(t *testing.T) {
	fn := exturi.URISetQuery()
	_, err := fn([]any{"https://example.com", "not-a-map"}, nil)
	if err == nil {
		t.Error("expected error for non-object params")
	}
}

// ---------- All ----------

func TestAll_Keys(t *testing.T) {
	m := exturi.All()
	expected := []string{"uriParse", "uriBuild", "uriJoin", "queryParse", "queryBuild", "uriGetPath", "uriGetQuery", "uriSetQuery"}
	for _, k := range expected {
		if _, ok := m[k]; !ok {
			t.Errorf("All(): missing key %q", k)
		}
	}
	if len(m) != len(expected) {
		t.Errorf("All(): expected %d keys, got %d", len(expected), len(m))
	}
}
