package helper

// TODO: Add comments

import (
	"log"
	"regexp"
)

func CompileRegex(expr *string) *regexp.Regexp {
	if expr == nil {
		return nil
	}
	regex, err := regexp.Compile(*expr)
	if err != nil {
		log.Fatalf("Failed to compile regex %s, error: %s", expr, err)
	}
	return regex
}
