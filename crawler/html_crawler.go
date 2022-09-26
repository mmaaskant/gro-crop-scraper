package crawler

import (
	"bytes"
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/attributes"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// HtmlCrawler crawls http(s) urls and returns their raw data,
// as it uses http.Client.Do() nothing is rendered so data hidden in API calls will not be fetched.
// HtmlCrawler is concurrency safe and keeps a registry of all found URLs.
type HtmlCrawler struct {
	*attributes.Tag
	client      *http.Client
	hrefRegex   *regexp.Regexp
	urlRegex    map[*regexp.Regexp]string
	urlRegistry map[string]string
	mutex       sync.RWMutex
}

func NewHtmlCrawler(c *http.Client) *HtmlCrawler {
	r, _ := regexp.Compile(`(href="(/?)((\w*)/)*")`)
	return &HtmlCrawler{
		nil,
		c,
		r,
		make(map[*regexp.Regexp]string),
		make(map[string]string),
		sync.RWMutex{},
	}
}

func (hc *HtmlCrawler) SetTag(t *attributes.Tag) {
	hc.Tag = t
}

// AddDiscoveryUrlRegex registers a new regex expression that is used to match URLs that should be collected for discovery.
func (hc *HtmlCrawler) AddDiscoveryUrlRegex(expr string) {
	hc.addRegex(expr, DiscoverRequestType)
}

// AddExtractUrlRegex registers a new regex expression that is used to match URLs that should be collected for extraction.
func (hc *HtmlCrawler) AddExtractUrlRegex(expr string) {
	hc.addRegex(expr, ExtractRequestType)
}

func (hc *HtmlCrawler) addRegex(expr string, requestType string) {
	r, err := regexp.Compile(expr)
	if err != nil {
		log.Fatalf("Failed to compile %s urlRegex %s, error: %s", requestType, expr, err)
	}
	hc.urlRegex[r] = requestType
}

// Crawl crawls the given Call and returns the data and URLs it has found while doing so.
func (hc *HtmlCrawler) Crawl(c *Call) *Data {
	body, err := hc.do(c.Request)
	if err != nil {
		log.Printf("Failed to crawl url: %s, error: %s", c.URL.String(), err)
		return NewData(hc.Tag, c, "", nil, err)
	}
	cleanedBody, err := hc.clean(c.Request, body)
	if err != nil {
		log.Printf("Failed to clean HTML fetched from url %s, error: %s. Skipping ...", c.Request.URL.String(), err)
	}
	calls := hc.findCalls(cleanedBody)
	return NewData(hc.Tag, c, body, calls, err)
}

// do calls the provided http.Request and cleans it by unescaping its body completing and any found partial urls.
func (hc *HtmlCrawler) do(req *http.Request) (string, error) {
	resp, err := hc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}

// findCalls uses the provided urlRegex to find urls and categorises them under either DiscoverRequestType or ExtractRequestType.
func (hc *HtmlCrawler) findCalls(body string) []*Call {
	calls := make([]*Call, 0)
	for regex, requestType := range hc.urlRegex {
		for _, url := range regex.FindAllString(body, -1) {
			if hc.hasRegisteredUrl(url) {
				hc.registerUrl(url, requestType)
				calls = append(calls, NewCall(NewRequest(http.MethodGet, url, nil), requestType))
			}
		}
	}
	return calls
}

// clean unescapes the given HTML it and completes any found partial urls.
func (hc *HtmlCrawler) clean(r *http.Request, body string) (string, error) {
	node, err := html.Parse(strings.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to parse HTML %s, error: %s", body, err)
		return "", err
	}
	err = hc.formatUrls(r, node)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = html.Render(buf, node)
	if err != nil {
		return "", err
	}
	body = html.UnescapeString(buf.String())
	body = strings.Replace(body, "\\", "", -1)
	body, err = hc.formatHiddenUrls(r, body)
	if err != nil {
		return "", err
	}
	return body, err
}

// formatUrls seeks out all href html attributes and parses them into full urls if required.
func (hc *HtmlCrawler) formatUrls(req *http.Request, n *html.Node) error {
	if n.Type == html.ElementNode {
		for k, attr := range n.Attr {
			if attr.Key == "href" && len(attr.Val) > 0 && attr.Val[0:1] != "#" {
				url, err := req.URL.Parse(attr.Val)
				if err != nil {
					return err
				}
				attr.Val = url.String()
				n.Attr[k] = attr
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if err := hc.formatUrls(req, child); err != nil {
			return err
		}
	}
	return nil
}

// formatHiddenUrls find any HTML href attributes that are not directly attached to any HTML tags,
// in places like scripts, other attributes, etc. These are then replaced with fully parsed URLs.
func (hc *HtmlCrawler) formatHiddenUrls(req *http.Request, body string) (string, error) {
	for _, href := range hc.hrefRegex.FindAllString(body, -1) {
		ref := strings.Trim(href, `href="`)
		url, err := req.URL.Parse(ref)
		if err != nil {
			return "", err
		}
		regex, err := regexp.Compile(href)
		if err != nil {
			return "", err
		}
		body = regex.ReplaceAllString(body, fmt.Sprintf(`href="%s"`, url.String()))
	}
	return body, nil
}

func (hc *HtmlCrawler) hasRegisteredUrl(url string) bool {
	hc.mutex.RLock()
	defer hc.mutex.RUnlock()
	_, ok := hc.urlRegistry[url]
	return !ok
}

func (hc *HtmlCrawler) registerUrl(url string, requestType string) {
	hc.mutex.Lock()
	defer hc.mutex.Unlock()
	hc.urlRegistry[url] = requestType
}
