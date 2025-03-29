package main

import (
	"fmt"

	"github.com/SimFG/etcd-mcp/tools"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

// Tool arguments are just structs, annotated with jsonschema tags
// More at https://mcpgolang.com/tools#schema-generation
type Content struct {
	Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
	Description *string `json:"description" jsonschema:"description=The description to submit"`
}
type MyFunctionsArguments struct {
	Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool (openai, google, claude, etc)"`
	Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(
		stdio.NewStdioServerTransport(),
		mcp_golang.WithName("Etcd MCP Server 🚀"),
		mcp_golang.WithVersion("1.0.0"),
	)
	err := server.RegisterTool("hello", "Say hello to a person", func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %server!", arguments.Submitter))), nil
	})
	if err != nil {
		panic(err)
	}

	tools.RegisterHealthTool(server)

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}
