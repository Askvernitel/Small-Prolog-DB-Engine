package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"weird/db/engine/client"
	"weird/db/engine/executor"
	"weird/db/engine/lexer"
	"weird/db/engine/parser"
)

type CLI struct {
	lexer    *lexer.Lexer
	executor *executor.Executor
}

func NewCLI(serverURL string) *CLI {
	if serverURL == "" {
		serverURL = "http://localhost:8081"
	}

	dbClient := client.NewClient(serverURL)

	return &CLI{
		lexer:    lexer.New(),
		executor: executor.NewExecutor(dbClient),
	}
}
func (c *CLI) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("╔═══════════════════════════════════════╗")
	fmt.Println("║   Welcome to WeirdDB SQL Shell        ║")
	fmt.Println("║   Connected to Prolog Backend         ║")
	fmt.Println("╚═══════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  - Type SQL statements (SELECT, INSERT, UPDATE, DELETE)")
	fmt.Println("  - 'exit' or 'quit' to exit")
	fmt.Println("  - 'help' for examples")
	fmt.Println()

	for {
		fmt.Print("weirddb> ")

		// Read input
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye! 👋")
			break
		}

		if input == "help" {
			c.printHelp()
			continue
		}

		if input == "" {
			continue
		}

		c.processStatement(input)

		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}

	c.executor.Close()
}

func (c *CLI) processStatement(input string) {
	tokens := c.lexer.Tokenize(input)

	if len(tokens) == 0 {
		fmt.Println("⚠️  No tokens generated")
		return
	}

	// Parse tokens
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
		return
	}

	if len(program.Statements) == 0 {
		fmt.Println("⚠️  No statements parsed")
		return
	}

	for _, stmt := range program.Statements {
		fmt.Printf("📝 Executing: %s\n", stmt.String())

		resp, err := c.executor.Execute(stmt)
		if err != nil {
			fmt.Printf("❌ Execution error: %v\n", err)
			continue
		}
		fmt.Println(resp.Message)

	}
}

func (c *CLI) printHelp() {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        SQL Examples                           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("SELECT Examples:")
	fmt.Println("  SELECT * FROM users")
	fmt.Println("  SELECT id, name, email FROM users")
	fmt.Println()
	fmt.Println("INSERT Examples:")
	fmt.Println("  INSERT INTO users (name, email, age) VALUES ('John', 'john@example.com', 30)")
	fmt.Println("  INSERT INTO products VALUES (1, 'Laptop', 999.99)")
	fmt.Println()
	fmt.Println("UPDATE Examples:")
	fmt.Println("  UPDATE users SET age = 31 WHERE name = 'John'")
	fmt.Println("  UPDATE products SET price = 899.99 WHERE id = 1")
	fmt.Println()
	fmt.Println("DELETE Examples:")
	fmt.Println("  DELETE FROM users WHERE name = 'John'")
	fmt.Println("  DELETE FROM products WHERE id = 1")
	fmt.Println()
}
