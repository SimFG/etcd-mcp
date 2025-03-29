package tools

import (
	"context"
	"errors"

	"github.com/SimFG/etcd-mcp/common"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type HealthArgument struct {
	ConnectionUrl string `json:"connection_url" jsonschema:"required,description=The connection URL for the etcd server"`
}

func RegisterHealthTool(server *mcp_golang.Server) {
	tooName := "health"
	toolDes := "Check the health of the etcd server"
	err := server.RegisterTool(tooName, toolDes, func(arguments HealthArgument) (*mcp_golang.ToolResponse, error) {
		err := common.EtcdOp(arguments.ConnectionUrl, func(ctx context.Context, client *clientv3.Client) error {
			_, err := client.Get(ctx, "health")
			if err == nil || errors.Is(err, rpctypes.ErrPermissionDenied) {
				return nil
			}
			return err
		})
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to check health: " + err.Error())), nil
		}
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("OK")), nil
	})
	if err != nil {
		panic(err)
	}
}
