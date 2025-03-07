package errors

import (
	"fmt"
	"strings"

	"github.com/ajtroup1/clear/object"
)

const (
	RED    = "\033[31m"
	YELLOW = "\033[33m"
	CLEAR  = "\033[0m"
)

type Error struct {
	Message   string
	Line      int
	Col       int
	Stage     string
	Context   string
	IsWarning bool
}

func New(message string, line, col int, stage string, lines []string, isWarning bool) *Error {
	context := lines[line-1]
	return &Error{Message: message, Line: line, Col: col, Stage: stage, Context: context, IsWarning: isWarning}
}

// Has errors checks if there are any errors in the lexer or parser
// The reason for two return values is that we need to handle differently if
// there are no errors but there are warnings and if there are just errors
// So the return is basically (hasErrors, hasWarnings)
func HasErrors(lexErrors, parseErrors []*Error) (bool, bool) {
	if len(lexErrors) > 0 || len(parseErrors) > 0 {
		for _, e := range lexErrors {
			if !e.IsWarning {
				return true, false
			}
		}
		for _, e := range parseErrors {
			if !e.IsWarning {
				return true, false
			}
		}
	}

	return false, len(lexErrors) > 0 || len(parseErrors) > 0
}

func ReportErrors(lexErrors, parseErrors []*Error) string {
	var errors string

	for _, e := range lexErrors {
		errors += report(e)
	}

	for _, e := range parseErrors {
		errors += report(e)
	}

	return errors
}

func ReportEvaluationError(err *object.Error) string {
	var out string
	out += RED + "Program evaluatation resulted in an error\n"
	out += fmt.Sprintf("\nEvaluation::Error [line: %d, col: %d] ---> %s.\n\tContext: %s\n"+CLEAR, err.Position.Line, err.Position.Col, Capitalize(err.Message), err.Context)
	return out
}

func report(e *Error) string {
	if !e.IsWarning {
		return fmt.Sprintf(RED+"%s::Error [line: %d, col: %d] ---> %s.\n\tContext: %s\n"+CLEAR, Capitalize(e.Stage), e.Line, e.Col, Capitalize(e.Message), e.Context)
	}

	return fmt.Sprintf(YELLOW+"%s::Warning [line: %d, col: %d] ---> %s.\n\tContext: %s\n"+CLEAR, Capitalize(e.Stage), e.Line, e.Col, Capitalize(e.Message), e.Context)
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
