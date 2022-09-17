package crawler

import (
	"bytes"
	"github.com/mmaaskant/gro-crop-scraper/attributes"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// HtmlCrawler crawls http(s) urls and returns their raw HTML,
// as it uses http.Get() nothing is rendered so data hidden in Javascript/APIs will not be fetched.
// HtmlCrawler is concurrency safe and keeps a registry of all found urls.
type HtmlCrawler struct {
	*attributes.Tag
	client      *http.Client
	regex       map[*regexp.Regexp]string
	urlRegistry map[string]string
	mutex       sync.RWMutex
}

// NewHtmlCrawler returns a new instance of HtmlCrawler and allows http.Client to be configured.
func NewHtmlCrawler(id string, c *http.Client) *HtmlCrawler {
	return &HtmlCrawler{
		attributes.NewTag("", id),
		c,
		make(map[*regexp.Regexp]string),
		make(map[string]string),
		sync.RWMutex{},
	}
}

// AddDiscoveryUrlRegex registers a new regex expression that is used to match URLs that should be collected for discovery.
func (hc *HtmlCrawler) AddDiscoveryUrlRegex(expr string) {
	hc.addRegex(expr, DiscoverRequestType)
}

// AddExtractUrlRegex registers a new regex expression that is used to match URLs that should be collected for extraction.
func (hc *HtmlCrawler) AddExtractUrlRegex(expr string) {
	hc.addRegex(expr, ExtractRequestType)
}

// addRegex registers a new regex expression in HtmlCrawler.
func (hc *HtmlCrawler) addRegex(expr string, t string) {
	r, err := regexp.Compile(expr)
	if err != nil {
		log.Fatalf("Failed to compile %s regex %s, error: %s", t, expr, err)
	}
	hc.regex[r] = t
}

// Crawl crawls the given call and returns the data and urls it has found while doing so.
func (hc *HtmlCrawler) Crawl(c *Call) *Data {
	b, err := hc.do(c.Request)
	if err != nil {
		log.Printf("Failed to crawl url: %s, error: %s", c.URL.String(), err)
		return NewCrawlerData(hc.Tag, c, "", nil, err)
	}
	calls := hc.findCalls(b)
	return NewCrawlerData(hc.Tag, c, b, calls, err)
}

// do calls the provided http.Request, cleans it by unescaping its body completing any found partial urls.
func (hc *HtmlCrawler) do(req *http.Request) (string, error) {
	resp, err := hc.client.Do(req)
	if err != nil {
		log.Printf("Failed to call url: %s, error: %s", req.URL.String(), err)
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to format body %v, error: %s", b, err)
		return "", err
	}
	return hc.clean(req, string(b))
}

// clean unescapes the given HTML it and completes any found partial urls.
func (hc *HtmlCrawler) clean(r *http.Request, b string) (string, error) {
	n, err := html.Parse(strings.NewReader(b))
	if err != nil {
		log.Printf("Failed to parse HTML %s, error: %s", b, err)
		return "", err
	}
	hc.formatUrls(r, n)
	buf := new(bytes.Buffer)
	err = html.Render(buf, n)
	if err != nil {
		log.Printf("Failed to render HTML %s, error: %s", b, err)
		return "", err
	}
	return html.UnescapeString(buf.String()), err
}

// findCalls uses the provided regex to find urls and categorises them under either DiscoverRequestType or ExtractRequestType.
func (hc *HtmlCrawler) findCalls(b string) []*Call {
	b = strings.Replace(b, "\\", "", -1)
	calls := make([]*Call, 0)
	for r, t := range hc.regex {
		for _, url := range r.FindAllString(b, -1) {
			if hc.hasRegisteredUrl(url) {
				hc.registerUrl(url, t)
				calls = append(calls, NewCrawlerCall(NewRequest(http.MethodGet, url, nil), t))
			}
		}
	}
	return calls
}

// hasRegisteredUrl checks if an url has been registered already and returns a bool accordingly.
func (hc *HtmlCrawler) hasRegisteredUrl(url string) bool {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()
	_, ok := hc.urlRegistry[url]
	return !ok
}

// registerUrl registers a new url in HtmlCrawler, preventing it from being visited again.
func (hc *HtmlCrawler) registerUrl(url string, t string) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.urlRegistry[url] = t
}

// formatUrls seeks out all href html attributes and parses them into full urls if required.
func (hc *HtmlCrawler) formatUrls(req *http.Request, n *html.Node) {
	if n.Type == html.ElementNode {
		for k, attr := range n.Attr {

			if attr.Key == "href" && len(attr.Val) > 0 && attr.Val[0:1] != "#" {
				url, err := req.URL.Parse(attr.Val)
				if err != nil {
					log.Fatalf("HTML Crawler failed to parse url %s, error: %s", attr.Val, err)
				}
				attr.Val = url.String()
				n.Attr[k] = attr
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		hc.formatUrls(req, child)
	}
}
