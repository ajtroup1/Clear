package com.jclear;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;
import java.util.Scanner;

public class JClear {
  static boolean hadError = false; // Global var keeping up with whether the program has encountered an error

  public static void main(String[] args) throws IOException {
    // JClear can be run in two ways:
    // - Executes a script when given a file path
    // - Usage: "jclear"
    // - Initiates a REPL instance
    // - Usage: "jclear [script].jc"
    if (args.length > 1) {
      System.out.println("Usage: jclear [script].jc");
      System.exit(64);
    } else if (args.length == 1) {
      // Argument detected, it must be a script to execute
      runFile(args[0]);
    } else {
      // No arguments, instantiate a REPL
      runPrompt();
    }
  }

  // Simply read the entire file's contents into a string and 'run's it
  private static void runFile(String path) throws IOException {
    byte[] bytes = Files.readAllBytes(Paths.get(path));
    run(new String(bytes, Charset.defaultCharset()));

    if (hadError)
      System.exit(65);
  }

  // Infinitely loop, reading input lines and 'run'ning them individually
  // REPL mode
  // Exit on NULL input
  private static void runPrompt() throws IOException {
    // Instantiate readers for continued user input
    InputStreamReader input = new InputStreamReader(System.in);
    BufferedReader reader = new BufferedReader(input);

    // Loop infinitely
    for (;;) {
      System.out.print(">> "); // Prompt indicator
      String line = reader.readLine(); // Read in the user's input
      // Exit the REPL on NULL input
      if (line == null)
        break;
      run(line); // Actually execute the input
      hadError = false;
    }
  }

  // Execute a literal string of source code
  private static void run(String source) {
    // Instantiate a scanner to tokenize the source string
    com.jclear.Scanner scanner = new com.jclear.Scanner(source);
    List<Token> tokens = scanner.scanTokens(); // Store the tokens into a list returned from scanTokens()

    // For now, just print the tokens.
    for (Token token : tokens) {
      System.out.println(token);
    }
  }

  // Reports an error message given line and context information
  public static void error(int line, String message) {
    report(line, "", message);
  }

  // Helper function to display the provided error info to the output
  private static void report(int line, String where,
      String message) {
    System.err.println(
        "[line " + line + "] Error" + where + ": " + message);
    hadError = true;
  }
}
