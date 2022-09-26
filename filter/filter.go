package filter

import "github.com/mmaaskant/gro-crop-scraper/attributes"

// TODO: Add comments

type Filter interface {
	attributes.Taggable
	Clone() Filter
	Filter(s string) map[string]string
}
