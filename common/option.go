package common

import "github.com/mark3labs/mcp-go/mcp"

func GetConnectionOptions() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("connection_url",
			mcp.DefaultString("http://localhost:2379"),
			mcp.Required(),
			mcp.Description("The connection URL for the etcd server"),
		),
	}
}
