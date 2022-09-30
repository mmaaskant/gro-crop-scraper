package helper

import (
	"log"
	"regexp"
)

// CompileRegex compiles the given regex expression and panics in case of an invalid format.
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
