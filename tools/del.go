package tools

import (
	"context"

	"github.com/SimFG/etcd-mcp/common"
	mcp_golang "github.com/metoro-io/mcp-golang"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type DelArgument struct {
	ConnectionUrl string `json:"connection_url" jsonschema:"required,description=The connection URL for the etcd server"`
	Key           string `json:"key" jsonschema:"required,description=The key to delete"`
}

func RegisterDelTool(server *mcp_golang.Server) {
	toolName := "del"
	toolDes := "Delete a key-value pair from etcd"
	err := server.RegisterTool(toolName, toolDes, func(arguments DelArgument) (*mcp_golang.ToolResponse, error) {
		err := common.EtcdOp(arguments.ConnectionUrl, func(ctx context.Context, client *clientv3.Client) error {
			_, err := client.Delete(ctx, arguments.Key)
			return err
		})
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to delete key: " + err.Error())), nil
		}
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("OK")), nil
	})
	if err != nil {
		panic(err)
	}
}
