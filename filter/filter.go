package filter

import "github.com/mmaaskant/gro-crop-scraper/attributes"

// TODO: Add comments

type Filter interface {
	attributes.Taggable
	Filter(data string) map[string]any
}
