package filter

import (
	"fmt"
	"github.com/mmaaskant/gro-crop-scraper/helper"
	"golang.org/x/net/html"
	"log"
	"reflect"
	"regexp"
)

// ConditionInterpreter functions as a Mediator for Condition, allowing the data to be typed and matched.
type ConditionInterpreter interface {
	Interpret(data any) bool
}

// KeyValueInterpreter implements ConditionInterpreter and allows a key value pair in the format
// map[string]any to be parsed.
type KeyValueInterpreter struct {
	condition *Condition
}

// Interpret implements ConditionInterpreter.Interpret.
func (kvi *KeyValueInterpreter) Interpret(data any) bool {
	var pair map[string]any
	pair, ok := data.(map[string]any)
	if !ok {
		log.Panicf(formatInterpreterTypeErrorMessage(pair, data))
	}
	for k, v := range pair {
		if !kvi.condition.MatchOne(&k, &v) {
			return false
		}
	}
	return true
}

func NewKeyValueInterpreter(keyExpr string, valueExpr string) *KeyValueInterpreter {
	return &KeyValueInterpreter{
		NewCondition(&keyExpr, &valueExpr),
	}
}

// HtmlTokenTagInterpreter implements ConditionInterpreter and allows an instance of *html.Token
// to be parsed, and checks if its tag matches or not.
type HtmlTokenTagInterpreter struct {
	condition *Condition
}

// Interpret implements ConditionInterpreter.Interpret.
func (htti *HtmlTokenTagInterpreter) Interpret(data any) bool {
	var token *html.Token
	token, ok := data.(*html.Token)
	if !ok {
		log.Panicf(formatInterpreterTypeErrorMessage(token, data))
	}
	return htti.condition.MatchOne(&token.Data, nil)
}

func NewHtmlTokenTagInterpreter(expr string) *HtmlTokenTagInterpreter {
	return &HtmlTokenTagInterpreter{
		NewCondition(&expr, nil),
	}
}

// HtmlTokenAttributeInterpreter implements ConditionInterpreter and allows an instance of *html.Token
// to be parsed, and checks if any matching attributes are found.
type HtmlTokenAttributeInterpreter struct {
	condition *Condition
}

// Interpret implements ConditionInterpreter.Interpret.
func (htai *HtmlTokenAttributeInterpreter) Interpret(data any) bool {
	var token *html.Token
	token, ok := data.(*html.Token)
	if !ok {
		log.Panicf(formatInterpreterTypeErrorMessage(token, data))
	}
	for _, attr := range token.Attr {
		if htai.condition.MatchOne(&attr.Key, &attr.Val) {
			return true
		}
	}
	return false
}

func NewHtmlTokenAttributeInterpreter(keyExpr string, valueExpr string) *HtmlTokenAttributeInterpreter {
	return &HtmlTokenAttributeInterpreter{
		NewCondition(&keyExpr, &valueExpr),
	}
}

func formatInterpreterTypeErrorMessage(expected any, got any) string {
	return fmt.Sprintf("Interperter expected type %s, got: %s", reflect.TypeOf(expected), reflect.TypeOf(got))
}

// Condition holds an optional key and value regex and determines if a value passes its requirements or not.
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

func (c *Condition) MatchOne(key *string, value any) bool {
	if c.keyRegex != nil && key != nil && !c.keyRegex.MatchString(*key) {
		return false
	}
	if c.valueRegex != nil && value != nil {
		switch v := value.(type) {
		case *string:
			if !c.valueRegex.MatchString(*v) {
				return false
			}
		case string:
			if !c.valueRegex.MatchString(v) {
				return false
			}
		case []any:
			for _, s := range v {
				if c.valueRegex.MatchString(fmt.Sprint(s)) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
	return true
}
