package main

import (
	"github.com/SimFG/etcd-mcp/tools"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(
		stdio.NewStdioServerTransport(),
		mcp_golang.WithName("Etcd MCP Server 🚀"),
		mcp_golang.WithVersion("1.0.0"),
	)

	tools.RegisterHealthTool(server)
	tools.RegisterGetTool(server)
	tools.RegisterPutTool(server)
	tools.RegisterDelTool(server)
	tools.RegisterStatusTool(server)
	tools.RegisterSaveTool(server)
	tools.RegisterRestoreTool(server)

	err := server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}
