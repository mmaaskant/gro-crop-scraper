package config

import (
	"flag"
	"fmt"
	"log"
	"strconv"
)

const (
	CrawlMethodStepId   string = "crawl"
	FilterMethodStepId  string = "filter"
	MapMethodStepId     string = "map"
	CompileMethodStepId string = "compile"
)

var stepFlags []*bool

// handleFlags registers flags based on the available steps and configs, and parses them.
// After the flags have been parsed, and if it includes step- and/or config flags,
// all steps and/or configs are omitted.
func handleFlags() {
	stepFlags = []*bool{
		flag.Bool(CrawlMethodStepId, false, fmt.Sprintf("Registers the %s step, which pulls raw data from external sources and saves it locally so it can be processed.", CrawlMethodStepId)),
		flag.Bool(FilterMethodStepId, false, fmt.Sprintf("Registers the %s step, which filters raw data and saves any data that is noteworthy.", FilterMethodStepId)),
		flag.Bool(MapMethodStepId, false, fmt.Sprintf("Registers the %s step, which maps filtered data to an universal format which makes it readable.", MapMethodStepId)),
		flag.Bool(CompileMethodStepId, false, fmt.Sprintf("Registers the %s step, which attempts to match mapped data based on their values.", CompileMethodStepId)),
	}
	for _, c := range configs {
		flag.Bool(c.Id, false, fmt.Sprintf("Registers the %s scraper config, runs all configs if none were registered.", c.Id))
	}
	flag.Parse()
	if hasFlaggedSteps() {
		applyFlaggedSteps()
	}
}

func hasFlaggedSteps() bool {
	for _, f := range stepFlags {
		if *f == true {
			return true
		}
	}
	return false
}

// applyFlaggedSteps omits any steps that have not been flagged only if any steps have been flagged.
func applyFlaggedSteps() {
	for _, c := range configs {
		for _, s := range c.Scrapers {
			if hasFlaggedSteps() {
				if !flagToBool(flag.Lookup(CrawlMethodStepId)) {
					s.Crawler = nil
				}
				if !flagToBool(flag.Lookup(FilterMethodStepId)) {
					s.Filter = nil
				}
				if !flagToBool(flag.Lookup(MapMethodStepId)) {
					//s.Mapper = nil
				}
				if !flagToBool(flag.Lookup(CompileMethodStepId)) {
					//s.Compiler = nil
				}
			}
		}
	}
}

func flagToBool(f *flag.Flag) bool {
	b, err := strconv.ParseBool(f.Value.String())
	if err != nil {
		log.Panicf("Failed to parse config flag into bool, error: %s", err)
	}
	return b
}
