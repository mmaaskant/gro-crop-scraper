package filter

import (
	"encoding/json"
	"log"
)

// JsonFilter implements Filter and iterates over the given JSON,
// As it walks over the given JSON it uses Criteria to search for any matching data.
// These matches are optionally extracted using Extractor and returned once the Filter has finished running.
type JsonFilter struct {
	*Tracker
}

func NewJsonFilter(criteria ...*Criteria) *JsonFilter {
	return &JsonFilter{
		NewFilterTracker(criteria),
	}
}

func (jf *JsonFilter) Clone() Filter {
	filterCopy := *jf
	trackerCopy := *jf.Tracker
	filterCopy.Tracker = &trackerCopy
	filterCopy.criteria = jf.criteria
	filterCopy.trackedCriteria = make(map[*Criteria]bool)
	return &filterCopy
}

func (jf *JsonFilter) Filter(s string) map[string]any {
	data := make(map[string]any, 0)
	var js map[string]any
	err := json.Unmarshal([]byte(s), &js)
	if err != nil {
		log.Printf("Failed to unmarshal JSON %s, error: %s", s, err)
		return nil
	}
	jf.Walk(js, data)
	return data
}

func (jf *JsonFilter) Walk(js map[string]any, data map[string]any, depths ...int) {
	var depth int
	if len(depths) == 0 {
		depths = append(depths, 1)
	}
	depth = depths[0]
	children := make([]*Criteria, 0)
	for k, v := range js {
		if walkable, ok := v.(map[string]any); ok {
			jf.Walk(walkable, data, depth+1)
		}
		for c, _ := range jf.getAllCriteria() {
			matched := c.Match(map[string]any{k: v})
			switch {
			case !matched && c.Parent == nil:
				delete(jf.trackedCriteria, c)
			case matched && c.Child != nil:
				children = append(children, c.Child)
				jf.trackedCriteria[c.Child] = false
			case matched && c.Child == nil && c.Extractor != nil:
				if extractedData := c.Extractor.Extract(map[string]any{k: v})[k]; extractedData != nil { // TODO: Doesn't support merge
					data[k] = extractedData
				}
			}
		}
	}
	for _, child := range children {
		delete(jf.trackedCriteria, child)
	}
}
