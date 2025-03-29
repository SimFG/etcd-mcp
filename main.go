package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/SimFG/etcd-mcp/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Etcd MCP Server 🚀",
		"1.0.0",
		server.WithLogging(),
	)

	// Add tool
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	// Add tool handler·
	s.AddTool(tool, helloHandler)

	tools.AddHealthTool(s)

	// // Start the stdio server
	// if err := server.ServeStdio(s); err != nil {
	// 	fmt.Printf("Server error: %v\n", err)
	// }

	sseServer := server.NewSSEServer(s, server.WithBaseURL("http://0.0.0.0:8181"))
	log.Printf("SSE server listening on :8181")
	if err := sseServer.Start(":8181"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
