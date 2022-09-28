package filter

import "github.com/mmaaskant/gro-crop-scraper/attributes"

// TODO: Add comments

type Filter interface {
	attributes.Taggable
	Clone() Filter
	Filter(s string) map[string]any
}

type Tracker struct {
	*attributes.Tag
	criteria        []*Criteria
	trackedCriteria map[*Criteria]bool
}

func NewFilterTracker(criteria []*Criteria) *Tracker {
	return &Tracker{
		nil,
		criteria,
		make(map[*Criteria]bool),
	}
}

func (tr *Tracker) SetTag(t *attributes.Tag) {
	tr.Tag = t
}

func (tr *Tracker) getAllCriteria() map[*Criteria]bool {
	for _, c := range tr.criteria {
		tr.trackedCriteria[c.Clone()] = false
	}
	return tr.trackedCriteria
}
