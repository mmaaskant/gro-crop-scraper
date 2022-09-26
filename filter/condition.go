package filter

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/helper"
	"golang.org/x/net/html"
	"log"
	"reflect"
	"regexp"
)

// TODO: Add comments

type ConditionInterpreter interface {
	Interpret(data any) bool
}

type HtmlTokenTagInterpreter struct {
	condition *Condition
}

func (htti *HtmlTokenTagInterpreter) Interpret(data any) bool {
	var token *html.Token
	token, ok := data.(*html.Token)
	if !ok {
		log.Fatalf(formatInterpreterTypeErrorMessage(token, data))
	}
	return htti.condition.MatchOne(&token.Data, nil)
}

func NewHtmlTokenTagInterpreter(expr string) *HtmlTokenTagInterpreter {
	return &HtmlTokenTagInterpreter{
		NewCondition(&expr, nil),
	}
}

type HtmlTokenAttributeInterpreter struct {
	condition *Condition
}

func (htai *HtmlTokenAttributeInterpreter) Interpret(data any) bool {
	var token *html.Token
	token, ok := data.(*html.Token)
	if !ok {
		log.Fatalf(formatInterpreterTypeErrorMessage(token, data))
	}
	for _, attr := range token.Attr {
		if !htai.condition.MatchOne(&attr.Key, &attr.Val) {
			return false
		}
	}
	return true
}

func NewHtmlTokenAttributeInterpreter(keyExpr string, valueExpr string) *HtmlTokenAttributeInterpreter {
	return &HtmlTokenAttributeInterpreter{
		NewCondition(&keyExpr, &valueExpr),
	}
}
func formatInterpreterTypeErrorMessage(expected any, got any) string {
	return fmt.Sprintf("Interperter expected type %s, got: %s", reflect.TypeOf(expected), reflect.TypeOf(got))
}

type Condition struct {
	keyRegex   *regexp.Regexp
	valueRegex *regexp.Regexp
}

func NewCondition(keyExpr *string, valueExpr *string) *Condition {
	return &Condition{
		helper.CompileRegex(keyExpr),
		helper.CompileRegex(valueExpr),
	}
}

func (c *Condition) MatchOne(key *string, value *string) bool {
	if c.keyRegex != nil && key != nil && !c.keyRegex.MatchString(*key) {
		return false
	}
	if c.valueRegex != nil && value != nil && !c.valueRegex.MatchString(*value) {
		return false
	}
	return true
}
