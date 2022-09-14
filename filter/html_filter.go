package filter

// TODO: Write html token filter + write one for JSON

type HtmlFilter struct {
}

func (hf *HtmlFilter) Filter() {
	//// TODO: Clean data? Only required if data is saved within struct
	//tz := html.NewTokenizer(strings.NewReader(s))
	//for tt := tz.Next(); tt != html.ErrorToken; tt = tz.Next() { // TODO: Also implement this for queues in Gophervisor package?
	//	t := tz.Token()
	//	switch tt {
	//	case html.SelfClosingTagToken, html.StartTagToken:
	//		hs.handleCriteria(&t)
	//	case html.TextToken:
	//		hs.scrapeTextByCriteria(t)
	//	case html.EndTagToken:
	//	case html.ErrorToken:
	//	}
	//}
	//return data
}

// TODO: 1 cleaner within html crawler, that uses nodes; this filter will use token instead and collect data, not required until extractor/collector step
