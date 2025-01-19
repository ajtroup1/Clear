package logger

import "os"

type Logger struct {
	out string
	// Since the lexer and parser run 'concurrently', we need to store the parser output separately
	// In this code, the parser lexes a token then decides whether to parse it
	// into a node or lex more tokens to parse a more complex node
	// It's not like the parser actually recieves a stream of tokens and analyzes them one after the other,
	// which is how the logs are structured to be easier to understand
	// So, parserOut will store output from the parser as it runs,
	// but actually dump it into the log after all lexing info is finished processing
	parserOut string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Append(s string) {
	l.out += s
}
func (l *Logger) AppendParser(s string) {
	l.parserOut += s
}

func (l *Logger) Get() string {
	return l.out
}
func (l *Logger) GetParserLog() string {
	return l.parserOut
}

func (l *Logger) WriteFile(filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.WriteString(l.out)
}

func (l *Logger) InitText(filepath string) {
	l.out += "# " + filepath + "\n"
	l.out += "Welcome to Clear\n\n*This file is a log of all activity that occured during the interpretation of your source code.*\n"
}

func (l *Logger) DefineSection(section, description string) {
	l.out += "\n"
	l.out += "## " + section
	l.out += "\n"
	l.out += description
	l.out += "\n\n"
}
