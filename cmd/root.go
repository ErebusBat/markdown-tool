package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/erebusbat/markdown-tool/internal/config"
	"github.com/erebusbat/markdown-tool/internal/parser"
	"github.com/erebusbat/markdown-tool/internal/writer"
	"github.com/erebusbat/markdown-tool/pkg/types"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "markdown-tool",
	Short: "Transform text inputs into well-formatted markdown",
	Long: `A lightweight command-line tool that processes text inputs and transforms them 
into well-formatted markdown suitable for knowledge management tools like Vimwiki and Obsidian.

The tool detects URLs (GitHub, JIRA, Notion, generic) and JIRA issue keys,
transforming them into appropriate markdown links.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			log.Fatalf("Error: %v", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/markdown-tool/config.yaml)")
}

func run() error {
	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get input from stdin or clipboard
	input, err := getInput()
	if err != nil {
		return fmt.Errorf("failed to get input: %w", err)
	}

	// Trim whitespace from input
	input = strings.TrimSpace(input)
	if input == "" {
		return nil // No input, nothing to do
	}

	// Parse input
	parsers := parser.GetParsers(cfg)
	contexts := make([]*types.ParseContext, 0)

	for _, p := range parsers {
		if ctx, err := p.Parse(input); err == nil && ctx != nil {
			contexts = append(contexts, ctx)
		}
	}

	// Vote on best writer
	writers := writer.GetWriters(cfg)
	bestWriter, bestScore := writer.Vote(writers, contexts)

	if bestWriter == nil || bestScore == 0 {
		// No writer wants to handle this, output verbatim
		fmt.Print(input)
		return nil
	}

	// Generate output
	if len(contexts) == 0 {
		fmt.Print(input)
		return nil
	}

	output, err := bestWriter.Write(contexts[0])
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Print(output)
	return nil
}

func getInput() (string, error) {
	// Check if we have stdin input
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped in
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(input), nil
	}

	// No stdin input, try clipboard
	return clipboard.ReadAll()
}
