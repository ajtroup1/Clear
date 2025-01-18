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

func New(message string, line, col int, stage, context string, isWarning bool) *Error {
	return &Error{Message: message, Line: line, Col: col, Stage: stage, Context: context, IsWarning: isWarning}
}

func HasErrors(lexErrors, parseErrors []Error) bool {
	if len(lexErrors) > 0 || len(parseErrors) > 0 {
		return true
	}

	return false
}

func ReportErrors(lexErrors, parseErrors []Error) string {
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
	out += RED + "Program evaluatation resulted in errors\n"
	out += fmt.Sprintf("\nEvaluation::Error ---> %s\n"+CLEAR, err.Message)
	return out
}

func report(e Error) string {
	if !e.IsWarning {
		return fmt.Sprintf(RED+"%s::Error [line: %d, col: %d] ---> %s\n\tcontext: '%s'\n"+CLEAR, capitalize(e.Stage), e.Line, e.Col, e.Message, e.Context)
	}

	return fmt.Sprintf(YELLOW+"%s::Warning [line: %d, col: %d] ---> %s\n\tcontext: '%s'\n"+CLEAR, capitalize(e.Stage), e.Line, e.Col, e.Message, e.Context)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
