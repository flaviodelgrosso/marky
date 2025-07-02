package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/flaviodelgrosso/marky"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Marky",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool
	tool := mcp.NewTool("convert_to_markdown",
		mcp.WithDescription("Convert a file to markdown format"),
		mcp.WithString("input",
			mcp.Required(),
			mcp.Description("Path to the input file to convert to markdown"),
		),
		mcp.WithString("output",
			mcp.Description("Path to the output markdown file"),
		),
	)

	// Add tool handler
	s.AddTool(tool, convertToMarkdown)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		log.Printf("Server error: %v\n", err)
	}
}

func convertToMarkdown(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	inputFile, err := request.RequireString("input")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	outputFile := request.GetString("output", "console")

	m := marky.New()
	result, err := m.Convert(inputFile)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to convert file: %v", err)), nil
	}

	if outputFile != "console" {
		if err := os.WriteFile(outputFile, []byte(result), 0o644); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write output file: %v", err)), nil
		}
	}

	return mcp.NewToolResultText(result), nil
}
