package helper

// TODO: Add comments

import (
	"log"
	"regexp"
)

func CompileRegex(expr *string) *regexp.Regexp {
	if expr == nil || *expr == "" {
		return nil
	}
	regex, err := regexp.Compile(*expr)
	if err != nil {
		log.Panicf("Failed to compile regex %s, error: %s", expr, err)
	}
	return regex
}
