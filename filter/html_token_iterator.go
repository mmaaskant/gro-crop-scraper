package filter

import (
	"golang.org/x/net/html"
	"strings"
)

// HtmlTokenIterator walks through a given HTML document using html.Tokenizer.
// The depth of the HTML document is also tracked and is available through Depth().
type HtmlTokenIterator struct {
	tokenizer *html.Tokenizer
	token     html.Token
	tags      []string
}

func newTokenIterator(s string) *HtmlTokenIterator {
	tz := html.NewTokenizer(strings.NewReader(s))
	return &HtmlTokenIterator{
		tz,
		tz.Token(),
		make([]string, 0),
	}
}

func (ti *HtmlTokenIterator) Next() html.TokenType {
	tokenType := ti.tokenizer.Next()
	ti.token = ti.tokenizer.Token()
	switch tokenType {
	case html.StartTagToken:
		ti.tags = append(ti.tags, ti.token.Data)
	case html.EndTagToken:
		for i := len(ti.tags) - 1; i >= 0; i-- {
			if ti.tags[i] == ti.token.Data {
				ti.tags = ti.tags[:i]
				break
			}
		}
	}
	return tokenType
}

func (ti *HtmlTokenIterator) Token() html.Token {
	return ti.token
}

func (ti *HtmlTokenIterator) Depth() int {
	return len(ti.tags)
}
