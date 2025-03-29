package tools

import (
	"context"
	"errors"

	"github.com/SimFG/etcd-mcp/common"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func AddHealthTool(s *server.MCPServer) {
	opts := []mcp.ToolOption{
		mcp.WithDescription("Check the health of the etcd server"),
	}
	opts = append(opts, common.GetConnectionOptions()...)
	tool := mcp.NewTool("health", opts...)
	s.AddTool(tool, healthHandler)
}

func healthHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	connectionURL := request.Params.Arguments["connection_url"].(string)
	err := common.EtcdOp(connectionURL, func(client *clientv3.Client) error {
		_, err := client.Get(ctx, "health")
		if err == nil || errors.Is(err, rpctypes.ErrPermissionDenied) {
			return nil
		}
		return err
	})
	if err != nil {
		return mcp.NewToolResultText("Failed to check health: " + err.Error()), nil
	}

	return mcp.NewToolResultText("OK"), nil
}
