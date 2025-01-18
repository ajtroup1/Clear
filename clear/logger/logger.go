package logger

import "os"

type Logger struct {
	out string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Append(s string) {
	l.out += s
}

func (l *Logger) Get() string {
	return l.out
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
	l.out += "# "+filepath+"\n"
	l.out += "Welcome to Clear\n\n*This file is a log of all activity that occured during the interpretation of your source code.*\n"
}

func (l *Logger) DefineSection(section, description string) {
	l.out += "\n"
	l.out += "## " + section
	l.out += "\n"
	l.out += description
	l.out += "\n"
}