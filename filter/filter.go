package filter

import "github.com/mmaaskant/gro-crop-scraper/attribute"

// Filter defines the ability to filter a type of data potentially extract any desired results using Criteria.
type Filter interface {
	attribute.Taggable
	Clone() Filter
	Filter(s string) map[string]any
}

// Tracker is used to track and manage a Filter's Criteria.
type Tracker struct {
	*attribute.Tag
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

func (tr *Tracker) SetTag(t *attribute.Tag) {
	tr.Tag = t
}

func (tr *Tracker) getAllCriteria() map[*Criteria]bool {
	for _, c := range tr.criteria {
		tr.trackedCriteria[c.Clone()] = false
	}
	return tr.trackedCriteria
}
