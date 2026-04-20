package exttext_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/exttext"
)

// ---------- WordCount ----------

func TestWordCount_Basic(t *testing.T) {
	fn := exttext.WordCount()
	got, err := fn([]any{"Hello world foo"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != float64(3) {
		t.Errorf("expected 3, got %v", got)
	}
}

func TestWordCount_Empty(t *testing.T) {
	fn := exttext.WordCount()
	got, _ := fn([]any{""}, nil)
	if got != float64(0) {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestWordCount_NoArgs(t *testing.T) {
	fn := exttext.WordCount()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestWordCount_WrongType(t *testing.T) {
	fn := exttext.WordCount()
	_, err := fn([]any{42}, nil)
	if err == nil {
		t.Error("expected error for non-string")
	}
}

// ---------- CharCount ----------

func TestCharCount_ASCII(t *testing.T) {
	fn := exttext.CharCount()
	got, _ := fn([]any{"hello"}, nil)
	if got != float64(5) {
		t.Errorf("expected 5, got %v", got)
	}
}

func TestCharCount_Unicode(t *testing.T) {
	fn := exttext.CharCount()
	got, _ := fn([]any{"héllo"}, nil)
	if got != float64(5) {
		t.Errorf("expected 5 runes, got %v", got)
	}
}

func TestCharCount_NoArgs(t *testing.T) {
	fn := exttext.CharCount()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestCharCount_WrongType(t *testing.T) {
	fn := exttext.CharCount()
	_, err := fn([]any{true}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- SentenceCount ----------

func TestSentenceCount_Basic(t *testing.T) {
	fn := exttext.SentenceCount()
	got, _ := fn([]any{"Hello. World! How are you?"}, nil)
	if got != float64(3) {
		t.Errorf("expected 3, got %v", got)
	}
}

func TestSentenceCount_NoPunctuation(t *testing.T) {
	fn := exttext.SentenceCount()
	got, _ := fn([]any{"Hello world"}, nil)
	if got != float64(1) {
		t.Errorf("expected 1, got %v", got)
	}
}

func TestSentenceCount_Empty(t *testing.T) {
	fn := exttext.SentenceCount()
	got, _ := fn([]any{""}, nil)
	if got != float64(0) {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestSentenceCount_NoArgs(t *testing.T) {
	fn := exttext.SentenceCount()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSentenceCount_WrongType(t *testing.T) {
	fn := exttext.SentenceCount()
	_, err := fn([]any{99}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- ReadingTime ----------

func TestReadingTime_Basic(t *testing.T) {
	fn := exttext.ReadingTime()
	// 200 words → 1 minute → 60 seconds
	words := make([]string, 200)
	for i := range words {
		words[i] = "word"
	}
	text := ""
	for i, w := range words {
		if i > 0 {
			text += " "
		}
		text += w
	}
	got, _ := fn([]any{text}, nil)
	if got != float64(60) {
		t.Errorf("expected 60s, got %v", got)
	}
}

func TestReadingTime_Empty(t *testing.T) {
	fn := exttext.ReadingTime()
	got, _ := fn([]any{""}, nil)
	if got != float64(0) {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestReadingTime_NoArgs(t *testing.T) {
	fn := exttext.ReadingTime()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestReadingTime_WrongType(t *testing.T) {
	fn := exttext.ReadingTime()
	_, err := fn([]any{false}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- WordFrequency ----------

func TestWordFrequency_Basic(t *testing.T) {
	fn := exttext.WordFrequency()
	got, err := fn([]any{"the cat and the dog"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	freq := got.(map[string]any)
	if freq["the"] != float64(2) {
		t.Errorf("expected the=2, got %v", freq["the"])
	}
	if freq["cat"] != float64(1) {
		t.Errorf("expected cat=1, got %v", freq["cat"])
	}
}

func TestWordFrequency_CaseInsensitive(t *testing.T) {
	fn := exttext.WordFrequency()
	got, _ := fn([]any{"Go go GO"}, nil)
	freq := got.(map[string]any)
	if freq["go"] != float64(3) {
		t.Errorf("expected go=3, got %v", freq["go"])
	}
}

func TestWordFrequency_Empty(t *testing.T) {
	fn := exttext.WordFrequency()
	got, _ := fn([]any{""}, nil)
	freq := got.(map[string]any)
	if len(freq) != 0 {
		t.Errorf("expected empty map, got %v", freq)
	}
}

func TestWordFrequency_NoArgs(t *testing.T) {
	fn := exttext.WordFrequency()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestWordFrequency_WrongType(t *testing.T) {
	fn := exttext.WordFrequency()
	_, err := fn([]any{3.14}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- Ngrams ----------

func TestNgrams_Bigrams(t *testing.T) {
	fn := exttext.Ngrams()
	got, err := fn([]any{"one two three four", float64(2)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	arr := got.([]any)
	if len(arr) != 3 {
		t.Errorf("expected 3 bigrams, got %d: %v", len(arr), arr)
	}
	if arr[0] != "one two" {
		t.Errorf("expected 'one two', got %v", arr[0])
	}
}

func TestNgrams_Unigrams(t *testing.T) {
	fn := exttext.Ngrams()
	got, _ := fn([]any{"a b c", float64(1)}, nil)
	arr := got.([]any)
	if len(arr) != 3 {
		t.Errorf("expected 3 unigrams, got %d", len(arr))
	}
}

func TestNgrams_NLargerThanWords(t *testing.T) {
	fn := exttext.Ngrams()
	got, _ := fn([]any{"hello", float64(5)}, nil)
	arr := got.([]any)
	if len(arr) != 0 {
		t.Errorf("expected empty, got %v", arr)
	}
}

func TestNgrams_NoArgs(t *testing.T) {
	fn := exttext.Ngrams()
	_, err := fn([]any{"hello"}, nil)
	if err == nil {
		t.Error("expected error for missing n")
	}
}

func TestNgrams_WrongFirstType(t *testing.T) {
	fn := exttext.Ngrams()
	_, err := fn([]any{42, float64(2)}, nil)
	if err == nil {
		t.Error("expected error for non-string s")
	}
}

func TestNgrams_InvalidN(t *testing.T) {
	fn := exttext.Ngrams()
	_, err := fn([]any{"hello world", float64(0)}, nil)
	if err == nil {
		t.Error("expected error for n=0")
	}
}

// ---------- Excerpt ----------

func TestExcerpt_Short(t *testing.T) {
	fn := exttext.Excerpt()
	got, err := fn([]any{"hello", float64(10)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Errorf("expected 'hello', got %v", got)
	}
}

func TestExcerpt_Truncated(t *testing.T) {
	fn := exttext.Excerpt()
	got, _ := fn([]any{"hello world", float64(5)}, nil)
	if got != "hello…" {
		t.Errorf("expected 'hello…', got %v", got)
	}
}

func TestExcerpt_ExactLength(t *testing.T) {
	fn := exttext.Excerpt()
	got, _ := fn([]any{"hello", float64(5)}, nil)
	if got != "hello" {
		t.Errorf("expected 'hello', got %v", got)
	}
}

func TestExcerpt_Unicode(t *testing.T) {
	fn := exttext.Excerpt()
	got, _ := fn([]any{"héllo world", float64(5)}, nil)
	if got != "héllo…" {
		t.Errorf("expected 'héllo…', got %v", got)
	}
}

func TestExcerpt_NoArgs(t *testing.T) {
	fn := exttext.Excerpt()
	_, err := fn([]any{"hello"}, nil)
	if err == nil {
		t.Error("expected error for missing maxLen")
	}
}

func TestExcerpt_WrongFirstType(t *testing.T) {
	fn := exttext.Excerpt()
	_, err := fn([]any{42, float64(5)}, nil)
	if err == nil {
		t.Error("expected error for non-string")
	}
}

func TestExcerpt_NegativeLen(t *testing.T) {
	fn := exttext.Excerpt()
	_, err := fn([]any{"hello", float64(-1)}, nil)
	if err == nil {
		t.Error("expected error for negative maxLen")
	}
}

// ---------- StripHTML ----------

func TestStripHTML_Basic(t *testing.T) {
	fn := exttext.StripHTML()
	got, err := fn([]any{"<b>Hello</b> <i>world</i>"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Hello world" {
		t.Errorf("expected 'Hello world', got %v", got)
	}
}

func TestStripHTML_NoTags(t *testing.T) {
	fn := exttext.StripHTML()
	got, _ := fn([]any{"plain text"}, nil)
	if got != "plain text" {
		t.Errorf("expected unchanged, got %v", got)
	}
}

func TestStripHTML_SelfClosing(t *testing.T) {
	fn := exttext.StripHTML()
	got, _ := fn([]any{"line1<br/>line2"}, nil)
	if got != "line1line2" {
		t.Errorf("expected 'line1line2', got %v", got)
	}
}

func TestStripHTML_NoArgs(t *testing.T) {
	fn := exttext.StripHTML()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestStripHTML_WrongType(t *testing.T) {
	fn := exttext.StripHTML()
	_, err := fn([]any{nil}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- NormalizeWhitespace ----------

func TestNormalizeWhitespace_Spaces(t *testing.T) {
	fn := exttext.NormalizeWhitespace()
	got, err := fn([]any{"hello   world"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %v", got)
	}
}

func TestNormalizeWhitespace_Tabs(t *testing.T) {
	fn := exttext.NormalizeWhitespace()
	got, _ := fn([]any{"hello\t\tworld"}, nil)
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %v", got)
	}
}

func TestNormalizeWhitespace_LeadingTrailing(t *testing.T) {
	fn := exttext.NormalizeWhitespace()
	got, _ := fn([]any{"  hello world  "}, nil)
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %v", got)
	}
}

func TestNormalizeWhitespace_Newlines(t *testing.T) {
	fn := exttext.NormalizeWhitespace()
	got, _ := fn([]any{"hello\n\nworld"}, nil)
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %v", got)
	}
}

func TestNormalizeWhitespace_NoArgs(t *testing.T) {
	fn := exttext.NormalizeWhitespace()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestNormalizeWhitespace_WrongType(t *testing.T) {
	fn := exttext.NormalizeWhitespace()
	_, err := fn([]any{123}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- All ----------

func TestAll_Keys(t *testing.T) {
	m := exttext.All()
	expected := []string{
		"wordCount", "charCount", "sentenceCount", "readingTime",
		"wordFrequency", "ngrams", "excerpt", "stripHTML", "normalizeWhitespace",
	}
	for _, k := range expected {
		if _, ok := m[k]; !ok {
			t.Errorf("All(): missing key %q", k)
		}
	}
	if len(m) != len(expected) {
		t.Errorf("All(): expected %d keys, got %d", len(expected), len(m))
	}
}
