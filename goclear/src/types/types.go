package types

type ParserError struct {
	Message     string
	Line        int
	Column      int
	LineContent string
}