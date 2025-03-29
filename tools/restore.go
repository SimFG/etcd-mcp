package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/SimFG/etcd-mcp/common"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"go.etcd.io/etcd/client/pkg/v3/logutil"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/etcdutl/v3/snapshot"
	"go.uber.org/zap"
)

type RestoreArgument struct {
	ConnectionUrl  string `json:"connection_url" jsonschema:"required,description=The connection URL for the etcd server"`
	FilePath       string `json:"file_path" jsonschema:"required,description=The file path to restore from"`
	IsSnapshotFile bool   `json:"is_snapshot_file" jsonschema:"required,description=Whether the restore file is a snapshot file"`
	OutDir         string `json:"out_dir" jsonschema:"description=Specify the etcd file directory when restoring snapshot files"`
}

func RegisterRestoreTool(server *mcp_golang.Server) {
	toolName := "restore"
	toolDes := "Restore key-value pairs from a file to etcd"
	err := server.RegisterTool(toolName, toolDes, func(arguments RestoreArgument) (*mcp_golang.ToolResponse, error) {
		if arguments.IsSnapshotFile {
			if arguments.OutDir == "" {
				return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Should give the out_dir param when restore a snapshot file"))), nil
			}
			lg, err := logutil.CreateDefaultZapLogger(zap.InfoLevel)
			if err != nil {
				return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Fail to create zap logger"))), nil
			}
			sp := snapshot.NewV3(lg)
			if err := sp.Restore(snapshot.RestoreConfig{
				SnapshotPath:  arguments.FilePath,
				OutputDataDir: arguments.OutDir,
				PeerURLs:      []string{arguments.ConnectionUrl},
			}); err != nil {
				return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to restore snapshot file: " + err.Error())), nil
			}
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Successfully restored the sanpshot file"))), nil
		}
		// Read the backup file
		content, err := os.ReadFile(arguments.FilePath)
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to read file: " + err.Error())), nil
		}

		var kvs []KeyValue
		if err := json.Unmarshal(content, &kvs); err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to unmarshal data: " + err.Error())), nil
		}

		err = common.EtcdOp(arguments.ConnectionUrl, func(ctx context.Context, client *clientv3.Client) error {
			for _, kv := range kvs {
				_, err := client.Put(ctx, kv.Key, kv.Value)
				if err != nil {
					return fmt.Errorf("failed to restore key %s: %v", kv.Key, err)
				}
			}
			return nil
		})
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to restore data: " + err.Error())), nil
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Successfully restored %d key-value pairs", len(kvs)))), nil
	})
	if err != nil {
		panic(err)
	}
}
