package errorlogger

import "fmt"

const (
	RED    = "\033[31m"
	YELLOW = "\033[33m"
	RESET  = "\033[0m"
)

type ErrorLogger struct {
	Lines  []string
	Debug  bool
	Errors []*Error
}

type Error struct {
	Line      int
	Col       int
	Message   string
	Context   string
	Stage     string
	IsWarning bool
}

func New(lines []string, debug bool) *ErrorLogger {
	return &ErrorLogger{
		Lines:  lines,
		Debug:  debug,
		Errors: make([]*Error, 0),
	}
}

func (el *ErrorLogger) SetLines(lines []string) {
	el.Lines = lines
}

func (el *ErrorLogger) NewError(line, col int, message, stage string, isWarning bool) {
	context := el.Lines[line-1]
	err := &Error{
		Line:      line,
		Col:       col,
		Message:   message,
		Context:   context,
		Stage:     stage,
		IsWarning: isWarning,
	}
	el.Errors = append(el.Errors, err)
}

func (el *ErrorLogger) ReportErrors() {
	for _, err := range el.Errors {
		if err.IsWarning {
			fmt.Printf(YELLOW+"%s::Warning: %s\n  [line: %d, col: %d]: '%s'\n", err.Stage, err.Message, err.Line, err.Col, el.Lines[err.Line-1])
		} else {
			fmt.Printf(RED+"%s::Error: %s\n  [line: %d, col: %d]: '%s'\n", err.Stage, err.Message, err.Line, err.Col, el.Lines[err.Line-1])
		}
		fmt.Print(RESET)
	}
}
