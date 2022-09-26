package filter

import (
	"github.com/mmaaskant/gro-crop-scraper/attributes"
	"golang.org/x/net/html"
	"log"
	"reflect"
	"strings"
)

// TODO: Add comments

type HtmlFilter struct {
	*attributes.Tag
	criteria        []*Criteria
	trackedCriteria map[*Criteria]any
}

func (hf *HtmlFilter) Clone() Filter {
	c := *hf
	c.trackedCriteria = make(map[*Criteria]any)
	return &c
}

func (hf *HtmlFilter) SetTag(t *attributes.Tag) {
	hf.Tag = t
}

func (hf *HtmlFilter) getAllCriteria() map[*Criteria]any {
	for _, c := range hf.criteria {
		hf.trackedCriteria[c.Clone()] = nil
	}
	return hf.trackedCriteria
}

func NewHtmlFilter(tag *attributes.Tag, criteria ...*Criteria) *HtmlFilter {
	return &HtmlFilter{
		tag,
		criteria,
		make(map[*Criteria]any),
	}
}

func (hf *HtmlFilter) Filter(s string) map[string]string {
	data := make(map[string]string, 0)
	ti := newTokenIterator(s)
	for tt := ti.Next(); tt != html.ErrorToken; tt = ti.Next() {
		t := ti.Token()
		for c, _ := range hf.getAllCriteria() {
			switch tt {
			case html.SelfClosingTagToken, html.StartTagToken:
				c.Depth = ti.Depth()
				if c.Match(&t) {
					switch {
					case c.Child != nil:
						hf.trackedCriteria[c.Child] = nil
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
					for c := range hf.trackedCriteria {
						if c.Child == nil && reflect.TypeOf(c.Extractor) == reflect.TypeOf((*HtmlTextExtractor)(nil)) {
							data = merge(data, c.Extractor.Extract(&t))
							delete(hf.trackedCriteria, c)
						}
					}
				}
			case html.EndTagToken:
				if c.Depth < ti.Depth() {
					if c.Parent != nil && c.Parent.Child != nil {
						hf.trackedCriteria[c.Parent] = nil
					}
					delete(hf.trackedCriteria, c)
				}
			}
		}
	}
	return data
}

func merge(destination map[string]string, source map[string]string) map[string]string { // TODO: Currently does not handle merging of data
	for k, v := range source {
		if _, hasKey := destination[k]; hasKey {
			log.Printf("Filter data aready has key %s with value %s, overwriting ...", k, destination[k])
		}
		destination[k] = v
	}
	return destination
}
