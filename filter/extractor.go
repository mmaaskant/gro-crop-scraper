package filter

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/helper"
	"golang.org/x/net/html"
	"log"
	"reflect"
	"regexp"
)

// TODO: Add comments

type Extractor interface {
	Extract(data any) map[string]string
}

type HtmlTextExtractor struct {
	id    string
	Clean *func(data map[string]string) map[string]string
}

func (hte *HtmlTextExtractor) Extract(data any) map[string]string {
	var token *html.Token
	token, ok := data.(*html.Token)
	if !ok {
		log.Fatalf(formatExtractorTypeErrorMessage(token, data))
	}
	text := map[string]string{hte.id: token.Data}
	if hte.Clean != nil {
		text = (*hte.Clean)(text)
	}
	return text
}

func NewHtmlTextExtractor(id string, clean ...func(map[string]string) map[string]string) *HtmlTextExtractor {
	var f *func(data map[string]string) map[string]string
	if len(clean) > 0 {
		f = &clean[0]
	}
	return &HtmlTextExtractor{
		id,
		f,
	}
}

type HtmlAttributeExtractor struct {
	keyRegex *regexp.Regexp
	Clean    *func(data map[string]string) map[string]string
}

func (hae *HtmlAttributeExtractor) Extract(data any) map[string]string {
	var token *html.Token
	token, ok := data.(*html.Token)
	if !ok {
		log.Fatalf(formatExtractorTypeErrorMessage(token, data))
	}
	attributes := make(map[string]string)
	for _, attr := range token.Attr {
		if hae.keyRegex.MatchString(attr.Key) {
			attributes[attr.Key] = attr.Val
		}
	}
	if hae.Clean != nil {
		attributes = (*hae.Clean)(attributes)
	}
	return attributes
}

func NewHtmlAttributeExtractor(keyExpr string, clean ...func(map[string]string) map[string]string) *HtmlAttributeExtractor {
	var f *func(data map[string]string) map[string]string
	if len(clean) > 0 {
		f = &clean[0]
	}
	return &HtmlAttributeExtractor{
		helper.CompileRegex(&keyExpr),
		f,
	}
}

func formatExtractorTypeErrorMessage(expected any, got any) string {
	return fmt.Sprintf("Extractor expected type %s, got: %s", reflect.TypeOf(expected), reflect.TypeOf(got)) // TODO: Implement this for all type of checks
}
