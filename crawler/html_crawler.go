package crawler

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	netUrl "net/url"
	"regexp"
	"strings"
	"sync"
)

// HtmlCrawler crawls http(s) urls and returns their raw HTML,
// as it uses http.Get() nothing is rendered so data hidden in Javascript/APIs will not be fetched.
// HtmlCrawler is concurrency safe and keeps a registry of all found urls.
type HtmlCrawler struct {
	tag         string
	origin      string
	client      *http.Client
	regex       map[*regexp.Regexp]string
	urlRegistry map[string]string
	mutex       sync.RWMutex
}

// NewHtmlCrawler returns a new instance of HtmlCrawler and allows http.Client to be configured.
func NewHtmlCrawler(origin string, client *http.Client) *HtmlCrawler {
	return &HtmlCrawler{
		"",
		origin,
		client,
		make(map[*regexp.Regexp]string),
		make(map[string]string),
		sync.RWMutex{},
	}
}

// GetTag returns HtmlCrawler's tag which is used to identify its data in other processes.
func (hc *HtmlCrawler) GetTag() string {
	return hc.tag
}

// SetTag sets HtmlCrawler's tag which is used to identify its data in other processes.
func (hc *HtmlCrawler) SetTag(tag string) {
	hc.tag = tag
}

// GetOrigin returns HtmlCrawler's origin which shows where the data originates from.
func (hc *HtmlCrawler) GetOrigin() string {
	return hc.tag
}

// AddDiscoveryUrlRegex registers a new regex expression that is used to match URLs that should be collected for discovery.
func (hc *HtmlCrawler) AddDiscoveryUrlRegex(expr string) {
	hc.addRegex(expr, DiscoverUrlType)
}

// AddExtractUrlRegex registers a new regex expression that is used to match URLs that should be collected for extraction.
func (hc *HtmlCrawler) AddExtractUrlRegex(expr string) {
	hc.addRegex(expr, ExtractUrlType)
}

// addRegex registers a new regex expression in HtmlCrawler
func (hc *HtmlCrawler) addRegex(expr string, t string) {
	r, err := regexp.Compile(expr)
	if err != nil {
		log.Fatalf("Failed to compile %s regex %s, error: %s", t, expr, err)
	}
	hc.regex[r] = t
}

// Crawl crawls the given call and returns the data and urls it has found while doing so.
func (hc *HtmlCrawler) Crawl(c *Call) *Data {
	b, err := hc.get(c.Url)
	if err != nil {
		log.Printf("Failed to crawl url: %s, error: %s", c.Url, err)
		return NewCrawlerData(hc.tag, hc.origin, c, "", nil, err)
	}
	calls := hc.findCalls(b)
	return NewCrawlerData(hc.tag, hc.origin, c, b, calls, err)
}

// get requests data from the given url, cleans it by unescapes it and completing any found partial urls.
func (hc *HtmlCrawler) get(url string) (string, error) {
	r, err := hc.client.Get(url)
	if err != nil {
		log.Printf("Failed to GET url: %s, error: %s", url, err)
		return "", err
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to format body %v, error: %s", b, err)
		return "", err
	}
	return hc.clean(url, string(b))
}

// clean unescapes the given HTML it and completes any found partial urls.
func (hc *HtmlCrawler) clean(url string, b string) (string, error) {
	n, err := html.Parse(strings.NewReader(b))
	if err != nil {
		log.Printf("Failed to parse HTML %s, error: %s", b, err)
		return "", err
	}
	hc.formatUrl(url, n)
	buf := new(bytes.Buffer)
	err = html.Render(buf, n)
	if err != nil {
		log.Printf("Failed to render HTML %s, error: %s", b, err)
		return "", err
	}
	return html.UnescapeString(buf.String()), err
}

// findCalls uses the provided regex to find urls and categorises them under either DiscoverUrlType or ExtractUrlType.
func (hc *HtmlCrawler) findCalls(b string) []*Call {
	b = strings.Replace(b, "\\", "", -1)
	urls := make([]*Call, 0)
	for r, t := range hc.regex {
		for _, url := range r.FindAllString(b, -1) {
			if hc.hasRegisteredUrl(url) {
				hc.registerUrl(url, t)
				urls = append(urls, NewCrawlerCall(url, t, http.MethodGet, nil, nil))
			}
		}
	}
	return urls
}

// registerUrl registers a new url in HtmlCrawler, preventing it from being visited again.
func (hc *HtmlCrawler) registerUrl(url string, t string) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.urlRegistry[url] = t
}

// hasRegisteredUrl checks if an url has been registered already and returns a bool accordingly.
func (hc *HtmlCrawler) hasRegisteredUrl(url string) bool {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()
	_, ok := hc.urlRegistry[url]
	return !ok
}

// formatUrl seeks out all href html attributes and if they are partial links will append
// either the current Call's base url or current url based on its format.
func (hc *HtmlCrawler) formatUrl(url string, n *html.Node) {
	url = hc.cleanUrl(url)
	if n.Type == html.ElementNode {
		for k, attr := range n.Attr {
			if attr.Key == "href" && attr.Val != "" {
				if ok, _ := hc.isCompleteUrl(attr.Val); !ok {
					if attr.Val[0:1] == "/" {
						attr.Val = hc.getBaseUrl(url) + attr.Val
					} else if attr.Val[0:1] != "#" {
						attr.Val = url + attr.Val
					}
					n.Attr[k] = attr
				}
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		hc.formatUrl(url, child)
	}
}

// getBaseUrl extracts the base url from the given url.
func (hc *HtmlCrawler) getBaseUrl(url string) string {
	u, _ := netUrl.Parse(url)
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

// cleanUrl makes sure that the provided url ends with "/" so it can be completed without corrupting.
func (hc *HtmlCrawler) cleanUrl(url string) string {
	lc := url[len(url)-1:]
	if lc != "/" {
		url = url + "/"
	}
	return url
}

// isCompleteUrl checks if the given url is callable or not and returns a bool accordingly.
func (hc *HtmlCrawler) isCompleteUrl(url string) (bool, error) {
	return regexp.MatchString(`^(https?://)(\S)*(\.[a-z]{2,5})`, url)
}
