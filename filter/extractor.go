package filter

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/helper"
	"golang.org/x/net/html"
	"log"
	"reflect"
	"regexp"
)

// Extractor attempts to map, organise and extract data and return it.
type Extractor interface {
	Extract(data any) map[string]any
}

// KeyValueExtractor implements Extractor and extracts data based on the provided key and value regexes.
// An optional Clean function can be provided which is run on the results found before they are returned.
type KeyValueExtractor struct {
	keyRegex   *regexp.Regexp
	valueRegex *regexp.Regexp
	Clean      *func(data map[string]any) map[string]any
}

func (kve *KeyValueExtractor) Extract(data any) map[string]any {
	var pair map[string]any
	pair, ok := data.(map[string]any)
	if !ok {
		log.Panicf(formatExtractorTypeErrorMessage(pair, data))
	}
	for k, v := range pair {
		if kve.keyRegex != nil {
			if !kve.keyRegex.MatchString(k) {
				return nil
			}
		}
		if kve.valueRegex != nil {
			switch value := v.(type) {
			case string:
				if !kve.valueRegex.MatchString(value) {
					return nil
				}
			case *string:
				if !kve.valueRegex.MatchString(*value) {
					return nil
				}
			case []any:
				values := make([]string, 0)
				for _, listItem := range value {
					item := fmt.Sprint(listItem)
					if kve.valueRegex.MatchString(item) {
						values = append(values, item)
					}
				}
				if len(values) == 0 {
					return nil
				}
				pair = map[string]any{k: values}
			}
		}

	}
	if kve.Clean != nil {
		pair = (*kve.Clean)(pair)
	}
	return pair
}

func NewKeyValueExtractor(keyExpr string, valueExpr string, clean ...func(data map[string]any) map[string]any) *KeyValueExtractor {
	var f *func(data map[string]any) map[string]any
	if len(clean) > 0 {
		f = &clean[0]
	}
	return &KeyValueExtractor{
		helper.CompileRegex(&keyExpr),
		helper.CompileRegex(&valueExpr),
		f,
	}
}

// HtmlTextExtractor implements Extractor and extracts text from the *html.Token.
// An optional Clean function can be provided which is run on the results found before they are returned.
type HtmlTextExtractor struct {
	id    string
	Clean *func(data map[string]any) map[string]any
}

func (hte *HtmlTextExtractor) Extract(data any) map[string]any {
	var token *html.Token
	token, ok := data.(*html.Token)
	if !ok {
		log.Panicf(formatExtractorTypeErrorMessage(token, data))
	}
	text := map[string]any{hte.id: token.Data}
	if hte.Clean != nil {
		text = (*hte.Clean)(text)
	}
	return text
}

func NewHtmlTextExtractor(id string, clean ...func(data map[string]any) map[string]any) *HtmlTextExtractor {
	var f *func(data map[string]any) map[string]any
	if len(clean) > 0 {
		f = &clean[0]
	}
	return &HtmlTextExtractor{
		id,
		f,
	}
}

// HtmlAttributeExtractor implements Extractor and extracts attributes from the given *html.Token.
// An optional Clean function can be provided which is run on the results found before they are returned.
type HtmlAttributeExtractor struct {
	keyRegex *regexp.Regexp
	Clean    *func(data map[string]any) map[string]any
}

func (hae *HtmlAttributeExtractor) Extract(data any) map[string]any {
	var token *html.Token
	token, ok := data.(*html.Token)
	if !ok {
		log.Panicf(formatExtractorTypeErrorMessage(token, data))
	}
	attributes := make(map[string]any)
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

func NewHtmlAttributeExtractor(keyExpr string, clean ...func(map[string]any) map[string]any) *HtmlAttributeExtractor {
	var f *func(data map[string]any) map[string]any
	if len(clean) > 0 {
		f = &clean[0]
	}
	return &HtmlAttributeExtractor{
		helper.CompileRegex(&keyExpr),
		f,
	}
}

func formatExtractorTypeErrorMessage(expected any, got any) string {
	return fmt.Sprintf("Extractor expected type %s, got: %s", reflect.TypeOf(expected), reflect.TypeOf(got))
}
