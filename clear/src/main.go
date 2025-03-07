// ? How can concurrency be used to improve the performance of Clear?

// * AI Summary:

// * Concurrency can be used to improve the performance of Clear by allowing multiple tasks to be executed
// * simultaneously. For example, when running scripts in a folder, each script can be executed in a separate
// * goroutine, which can significantly reduce the total execution time. Additionally, concurrency can be used
// * to parallelize the execution of independent tasks within a script, such as evaluating expressions or
// * executing statements. By leveraging goroutines and channels, Clear can take advantage of the inherent
// * parallelism in many programming tasks, leading to faster execution and improved performance.


package main

import (
	"flag"
	"fmt"
	"strings"

	helpPkg "github.com/ajtroup1/clear/src/help"
	"github.com/ajtroup1/clear/src/repl"
)

func main() {
	// Parse the command line arguments
	/*
	   Clear can be run in two modes: REPL and script.
	   If no arguments are provided, start the REPL.
	   If a file is provided, run the script.
	   If a folder is provided, run all scripts in the folder.
	   This is indicated by the path argument ending in a `/` or `\`.
	   If the folder contains a subfolder named "tests", run all tests in that folder.
	   Flags:
	   Ran with the command `clear [path] [flag(s)]`:
	   - [-d, --debug]: Enable debug mode
	   - [-t, --test]: Run all tests in the folder
	   Ran with the command `clear [flag(s)]`:
	   - [-v, --version]: Print the version of Clear
	   - [-h, --help]: Print the help message
	*/

	// Define and evaluate flags
	var debug, test, version, help bool

	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.BoolVar(&debug, "d", false, "Enable debug mode (short)")
	flag.BoolVar(&test, "test", false, "Run all tests in the folder")
	flag.BoolVar(&test, "t", false, "Run all tests in the folder (short)")
	flag.BoolVar(&version, "version", false, "Print the version of Clear")
	flag.BoolVar(&version, "v", false, "Print the version of Clear (short)")
	flag.BoolVar(&help, "help", false, "Print the help message")
	flag.BoolVar(&help, "h", false, "Print the help message (short)")
	flag.Parse()

	args := flag.Args()

	// Handle flags with no path argument
	if version {
		fmt.Println("Clear version 1.0.0")
		return
	}

	if help {
		helpPkg.PrintHelpText()
		return
	}

	// Run the REPL if no path argument is provided
	if len(args) == 0 {
		repl.StartREPL()
		return
	}

	// If the REPL was not run, there must be an argument which should be a path
	path := args[0]

	// Determine if the path is a script or a folder
	// If the path ends in a `/` or `\`, it is a folder
	// Otherwise, it is a single script
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		runFolder(path, test, debug)
	} else {
		runScript(path, debug)
	}
}

// runScript runs a single script.
// path: The path to the script.
// debug: Whether or not debug mode is enabled.
func runScript(path string, debug bool) {
	fmt.Printf("Running script: %s\n", path)
}

// runFolder runs all scripts in a directory recursively.
// path: The path to the folder.
// test: Whether or not to run tests on the folder.
// debug: Whether or not debug mode is enabled.
func runFolder(path string, test, debug bool) {
	fmt.Printf("Running all scripts in folder: %s\n", path)
	if test {
		fmt.Println("Running tests...")
	} else {
		// Folder execution logic here
	}
}
