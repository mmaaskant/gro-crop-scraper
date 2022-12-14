package scraper

import (
	"github.com/mmaaskant/gro-crop-scraper/attribute"
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"github.com/mmaaskant/gro-crop-scraper/filter"
)

// Scraper holds all components and implements attribute.Taggable,
// these components are used to execute their respective steps if they are available.
type Scraper struct {
	*attribute.Tag
	Crawler crawler.Crawler
	Calls   []*crawler.Call
	Filter  filter.Filter
	//Mapper mapper.Mapper
	//Compiler compiler.Compiler
}

func NewScraper(c crawler.Crawler, calls []*crawler.Call, f filter.Filter) *Scraper {
	return &Scraper{
		nil,
		c,
		calls,
		f,
	}
}

func (s *Scraper) SetTag(t *attribute.Tag) {
	s.Tag = t
	if s.Crawler != nil {
		s.Crawler.SetTag(t)
	}
	if s.Filter != nil {
		s.Filter.SetTag(t)
	}
}
