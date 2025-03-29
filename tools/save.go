package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/SimFG/etcd-mcp/common"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"go.etcd.io/etcd/client/pkg/v3/logutil"
	clientv3 "go.etcd.io/etcd/client/v3"
	snapshot "go.etcd.io/etcd/client/v3/snapshot"
	"go.uber.org/zap"
)

type SaveArgument struct {
	ConnectionUrl string `json:"connection_url" jsonschema:"required,description=The connection URL for the etcd server"`
	FilePath      string `json:"file_path" jsonschema:"required,description=The file path to save the backup"`
	StartKey      string `json:"start_key" jsonschema:"description=The start key to save the backup"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func RegisterSaveTool(server *mcp_golang.Server) {
	toolName := "save"
	toolDes := "Save all key-value pairs from etcd to a file"
	err := server.RegisterTool(toolName, toolDes, func(arguments SaveArgument) (*mcp_golang.ToolResponse, error) {
		var kvs []KeyValue
		err := common.EtcdOp(arguments.ConnectionUrl, func(ctx context.Context, client *clientv3.Client) error {
			if arguments.StartKey != "" {
				resp, err := client.Get(ctx, arguments.StartKey, clientv3.WithPrefix())
				if err != nil {
					return err
				}

				for _, kv := range resp.Kvs {
					kvs = append(kvs, KeyValue{
						Key:   string(kv.Key),
						Value: string(kv.Value),
					})
				}
				return nil
			}
			config := clientv3.Config{
				Endpoints: []string{arguments.ConnectionUrl},
				Context:   ctx,
			}
			lg, err := logutil.CreateDefaultZapLogger(zap.InfoLevel)
			if err != nil {
				return err
			}
			err = snapshot.Save(ctx, lg, config, arguments.FilePath)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to get key-values: " + err.Error())), nil
		}

		if arguments.StartKey == "" {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Successfully saved a snapshot file for all key-values")), nil
		}

		// Create directory if it doesn't exist
		dir := filepath.Dir(arguments.FilePath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to create directory: " + err.Error())), nil
		}

		// Save to file
		content, err := json.MarshalIndent(kvs, "", "  ")
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to marshal data: " + err.Error())), nil
		}

		if err := os.WriteFile(arguments.FilePath, content, 0o644); err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to write file: " + err.Error())), nil
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Successfully saved " + string(len(kvs)) + " key-value pairs")), nil
	})
	if err != nil {
		panic(err)
	}
}
