package config

import (
	"github.com/mmaaskant/gro-crop-scraper/crawler"
	"github.com/mmaaskant/gro-crop-scraper/filter"
	"github.com/mmaaskant/gro-crop-scraper/scraper"
	"log"
	"net/http"
	"reflect"
	"time"
)

const (
	BurpeeConfigId      = "burpee"
	BurpeeHtmlScraperId = "burpee_html"
)

// NewBurpeeConfig holds components configured to scrape https://burpee.com.
func NewBurpeeConfig() *Config {
	c := newConfig(BurpeeConfigId)
	c.AddScraper(BurpeeHtmlScraperId, scraper.NewScraper(newBurpeeHtmlCrawler(), getBurpeeHtmlCrawlerCalls(), getBurpeeHtmlFilter()))
	return c
}

// newBurpeeHtmlCrawler returns an instance of crawler.HtmlCrawler configured to scrape
// the Burpee website and a slice of crawler.Call instances to kick off the crawling process.
func newBurpeeHtmlCrawler() *crawler.HtmlCrawler {
	cr := crawler.NewHtmlCrawler(&http.Client{Timeout: 90 * time.Second})
	cr.AddDiscoveryUrlRegex(`(https?:\/\/)?www\.burpee\.com\/?(vegetables|flowers|perennials|herbs|fruit)([\w\/-]*)(\?p=\d{1,3})?(&is_scroll=1)?`)
	cr.AddExtractUrlRegex(`(https?:\/\/)?www\.burpee\.com\/([\w\-]*)(prod\d*.html)(\/)?`)
	return cr
}

func getBurpeeHtmlCrawlerCalls() []*crawler.Call {
	return []*crawler.Call{crawler.NewCall(
		crawler.NewRequest(http.MethodGet, "https://www.burpee.com", nil),
		crawler.DiscoverRequestType,
	)}
}

func getBurpeeHtmlFilter() *filter.HtmlFilter {
	cb := filter.NewCriteriaBuilder(
		filter.NewCriteria(
			nil,
			filter.NewHtmlTokenTagInterpreter("div"),
			filter.NewHtmlTokenAttributeInterpreter("class", "product-add-form"),
		),
	)
	cb.AddChild(filter.NewCriteria(
		filter.NewHtmlTextExtractor("attributes", func(data map[string]any) map[string]any {
			attributes, ok := data["attributes"].(string)
			if !ok {
				log.Printf("Burpee HTML filter expected %s, got %s", reflect.TypeOf(attributes), reflect.TypeOf(data["attributes"]))
			}
			f := filter.NewJsonFilter(
				filter.NewCriteria(
					filter.NewKeyValueExtractor("", ""),
					filter.NewKeyValueInterpreter("(bp_).*|name|description|short_description", ""),
				),
			)
			return f.Filter(attributes)
		}),
		filter.NewHtmlTokenTagInterpreter("script"),
		filter.NewHtmlTokenAttributeInterpreter("type", "text/x-magento-init"),
	))
	return filter.NewHtmlFilter(cb.Build())
}
