package tools

import (
	"context"
	"encoding/json"

	"github.com/SimFG/etcd-mcp/common"
	mcp_golang "github.com/metoro-io/mcp-golang"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type StatusArgument struct {
	ConnectionUrl string `json:"connection_url" jsonschema:"required,description=The connection URL for the etcd server"`
}

type StatusResponse struct {
	Version   string `json:"version"`
	DbSize    int64  `json:"db_size"`
	Leader    uint64 `json:"leader"`
	RaftIndex uint64 `json:"raft_index"`
	RaftTerm  uint64 `json:"raft_term"`
	IsLearner bool   `json:"is_learner"`
	Errors    string `json:"errors,omitempty"`
}

func RegisterStatusTool(server *mcp_golang.Server) {
	toolName := "status"
	toolDes := "Get the status of the etcd server"
	err := server.RegisterTool(toolName, toolDes, func(arguments StatusArgument) (*mcp_golang.ToolResponse, error) {
		var response StatusResponse
		err := common.EtcdOp(arguments.ConnectionUrl, func(ctx context.Context, client *clientv3.Client) error {
			status, err := client.Status(ctx, arguments.ConnectionUrl)
			if err != nil {
				return err
			}

			response = StatusResponse{
				Version:   status.Version,
				DbSize:    status.DbSize,
				Leader:    status.Leader,
				RaftIndex: status.RaftIndex,
				RaftTerm:  status.RaftTerm,
				IsLearner: status.IsLearner,
			}
			return nil
		})
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to get status: " + err.Error())), nil
		}

		content, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to marshal response: " + err.Error())), nil
		}
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(content))), nil
	})
	if err != nil {
		panic(err)
	}
}
