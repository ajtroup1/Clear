package util

import "fmt"

func PrintErrorPanic(step string, msg string) {
	panic(fmt.Sprintf("\033[31m%s::Error -> %s\n\033[0m", step, msg))
}
