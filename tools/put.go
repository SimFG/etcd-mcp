package tools

import (
	"context"

	"github.com/SimFG/etcd-mcp/common"
	mcp_golang "github.com/metoro-io/mcp-golang"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type PutArgument struct {
	ConnectionUrl string `json:"connection_url" jsonschema:"required,description=The connection URL for the etcd server"`
	Key           string `json:"key" jsonschema:"required,description=The key to set"`
	Value         string `json:"value" jsonschema:"required,description=The value to set"`
}

func RegisterPutTool(server *mcp_golang.Server) {
	toolName := "put"
	toolDes := "Put a key-value pair into etcd"
	err := server.RegisterTool(toolName, toolDes, func(arguments PutArgument) (*mcp_golang.ToolResponse, error) {
		err := common.EtcdOp(arguments.ConnectionUrl, func(ctx context.Context, client *clientv3.Client) error {
			_, err := client.Put(ctx, arguments.Key, arguments.Value)
			return err
		})
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to put key-value: " + err.Error())), nil
		}
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("OK")), nil
	})
	if err != nil {
		panic(err)
	}
}
