package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajtroup1/goclear/lexer"
	"github.com/ajtroup1/goclear/parser"
	"github.com/sanity-io/litter"
)

func main() {
	jsonMode := true
	litterMode := false
	readingDir := false
	debug := false

	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path")
		return
	}

	filePath := os.Args[1]

	if err := clearAndCreateJsonFolder(); err != nil {
		fmt.Println("Error setting up /jsons folder:", err)
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("Error checking path:", err)
		return
	}

	if info.IsDir() {
		readingDir = true
		fmt.Println("Reading directory:", filePath)
	}

	if readingDir {
		files, err := os.ReadDir(filePath)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".clr" {
				processFile(filepath.Join(filePath, file.Name()), debug, jsonMode, litterMode)
			}
		}
	} else {
		processFile(filePath, debug, jsonMode, litterMode)
	}
}

func clearAndCreateJsonFolder() error {
	jsonsDir := "./jsons"

	if err := os.RemoveAll(jsonsDir); err != nil {
		return fmt.Errorf("failed to remove /jsons folder: %w", err)
	}

	if err := os.Mkdir(jsonsDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create /jsons folder: %w", err)
	}

	return nil
}

func processFile(filePath string, debug, jsonMode, litterMode bool) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if len(os.Args) > 2 && os.Args[2] == "-d" {
		debug = true
	}

	src := string(bytes)
	lexer := lexer.New(src)
	lexer.Lex()
	if debug {
		// for i, token := range lexer.Tokens {
		// 	fmt.Println(i, token.Stringify())
		// }
	}

	parser := parser.New(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) != 0 {
		fmt.Printf("\033[31mParser errors for '%s':\n", filePath)
		for _, err := range parser.Errors() {
			fmt.Printf("\tParser::Error --> '%s' [line: %d, col: %d]\n", err.Msg, err.Line, err.Col)
		}
		fmt.Print("\033[0m")
		return
	}

	if debug {
		if len(program.Statements) == 0 {
			fmt.Println("No program statements")
			return
		}
		
		if jsonMode {
			jsonFilePath := filepath.Join("jsons", generateJsonFilename(filePath))
			file, err := os.Create(jsonFilePath)
			if err != nil {
				fmt.Println("Error creating JSON file:", err)
				return
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(program); err != nil {
				fmt.Println("Error encoding program to JSON:", err)
				return
			}
			fmt.Printf("Program written to '%s'\n", jsonFilePath)
		}
		
		if litterMode {
			fmt.Println("Program Statements:")
			litter.Dump(program)
		}
	}
}

func generateJsonFilename(filePath string) string {
	baseName := filepath.Base(filePath)
	nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	return nameWithoutExt + ".json"
}