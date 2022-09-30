package filter

import (
	"golang.org/x/net/html"
	"log"
	"reflect"
	"strings"
)

// HtmlFilter implements Filter and iterates over HTML using HtmlTokenIterator.
// As it walks through the HTML document it searches for any matching Criteria,
// and extracts any found data with if Criteria has an Extractor.
type HtmlFilter struct {
	*Tracker
}

func NewHtmlFilter(criteria ...*Criteria) *HtmlFilter {
	return &HtmlFilter{
		NewFilterTracker(criteria),
	}
}

func (hf *HtmlFilter) Clone() Filter {
	filterCopy := *hf
	trackerCopy := *hf.Tracker
	filterCopy.Tracker = &trackerCopy
	filterCopy.criteria = hf.criteria
	filterCopy.trackedCriteria = make(map[*Criteria]bool)
	return &filterCopy
}

// Filter iterates over all tags within the given HTML, and applies Criteria for every found start tag.
// Any fully matched Criteria that have an Extractor will extract data from the matched tag
// and return it once the filter is finished.
func (hf *HtmlFilter) Filter(s string) map[string]any {
	data := make(map[string]any, 0)
	ti := newTokenIterator(s)
	for tt := ti.Next(); tt != html.ErrorToken; tt = ti.Next() {
		t := ti.Token()
		for c, _ := range hf.getAllCriteria() {
			switch tt {
			case html.SelfClosingTagToken, html.StartTagToken:
				c.Depth = ti.Depth()
				if c.Match(&t) {
					hf.trackedCriteria[c] = true
					switch {
					case c.Child != nil:
						hf.trackedCriteria[c.Child] = false
					case c.Child == nil:
						if reflect.TypeOf(c.Extractor) == reflect.TypeOf((*HtmlAttributeExtractor)(nil)) {
							data = merge(data, c.Extractor.Extract(&t))
						}
					}
				}
				if reflect.TypeOf(c.Extractor) != reflect.TypeOf((*HtmlTextExtractor)(nil)) {
					delete(hf.trackedCriteria, c)
				}
			case html.TextToken:
				if len(strings.TrimSpace(t.Data)) != 0 {
					for c, passed := range hf.trackedCriteria {
						if passed && c.Child == nil && reflect.TypeOf(c.Extractor) == reflect.TypeOf((*HtmlTextExtractor)(nil)) {
							if extractedData := c.Extractor.Extract(&t); extractedData != nil {
								data = merge(data, c.Extractor.Extract(&t))
							}
							delete(hf.trackedCriteria, c)
						}
					}
				}
			case html.EndTagToken:
				if c.Depth < ti.Depth() {
					if c.Parent != nil && c.Parent.Child != nil {
						hf.trackedCriteria[c.Parent] = false
					}
					delete(hf.trackedCriteria, c)
				}
			}
		}
	}
	return data
}

func merge(destination map[string]any, source map[string]any) map[string]any { // TODO: Currently does not handle merging of data
	for k, v := range source {
		if _, hasKey := destination[k]; hasKey {
			log.Printf("Filter data aready has key %s with value %s, overwriting ...", k, destination[k])
		}
		destination[k] = v
	}
	return destination
}
