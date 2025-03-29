package tools

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/SimFG/etcd-mcp/common"
	mcp_golang "github.com/metoro-io/mcp-golang"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type GetArgument struct {
	ConnectionUrl string `json:"connection_url" jsonschema:"required,description=The connection URL for the etcd server"`
	Key           string `json:"key" jsonschema:"required,description=The key to get"`
	WithPrefix    bool   `json:"with_prefix,omitempty" jsonschema:"description=Get all keys with the given prefix"`
	WithRange     string `json:"with_range,omitempty" jsonschema:"description=Get all keys in the range [key, range_end)"`
	WithOrder     string `json:"with_order,omitempty" jsonschema:"enum=asc,desc,description=Order of the results (asc or desc)"`
	KeyOnly       bool   `json:"key_only,omitempty" jsonschema:"description=If true, only return keys without values"`
	Limit         int    `json:"limit,omitempty" jsonschema:"description=Maximum number of keys to return"`
}

type GetResponse struct {
	Kvs      []KeyValuePair `json:"kvs"`
	Count    int64          `json:"count"`
	More     bool           `json:"more"`
	Revision int64          `json:"revision"`
}

type KeyValuePair struct {
	Key            string `json:"key"`
	Value          string `json:"value,omitempty"`
	CreateRevision int64  `json:"create_revision,omitempty"`
	ModRevision    int64  `json:"mod_revision,omitempty"`
	Version        int64  `json:"version,omitempty"`
}

func RegisterGetTool(server *mcp_golang.Server) {
	toolName := "get"
	toolDes := "Get key-value pairs from etcd with various options"
	err := server.RegisterTool(toolName, toolDes, func(arguments GetArgument) (*mcp_golang.ToolResponse, error) {
		var response GetResponse
		err := common.EtcdOp(arguments.ConnectionUrl, func(ctx context.Context, client *clientv3.Client) error {
			opts := []clientv3.OpOption{}

			// Handle WithPrefix
			if arguments.WithPrefix {
				opts = append(opts, clientv3.WithPrefix())
			}

			// Handle WithRange
			if arguments.WithRange != "" {
				opts = append(opts, clientv3.WithRange(arguments.WithRange))
			}

			// Handle WithOrder
			if arguments.WithOrder != "" {
				order := strings.ToLower(arguments.WithOrder)
				if order == "desc" {
					opts = append(opts, clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
				} else if order == "asc" {
					opts = append(opts, clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
				}
			}

			// Handle KeyOnly
			if arguments.KeyOnly {
				opts = append(opts, clientv3.WithKeysOnly())
			}

			// Handle Limit
			if arguments.Limit > 0 {
				opts = append(opts, clientv3.WithLimit(int64(arguments.Limit)))
			}

			// Add WithCountOnly to get total count
			getResp, err := client.Get(ctx, arguments.Key, append(opts, clientv3.WithCountOnly())...)
			if err != nil {
				return err
			}
			response.Count = getResp.Count

			// Get actual key-values
			resp, err := client.Get(ctx, arguments.Key, opts...)
			if err != nil {
				return err
			}

			response.Revision = resp.Header.Revision
			response.More = resp.More
			response.Kvs = make([]KeyValuePair, 0, len(resp.Kvs))

			for _, kv := range resp.Kvs {
				pair := KeyValuePair{
					Key:            string(kv.Key),
					CreateRevision: kv.CreateRevision,
					ModRevision:    kv.ModRevision,
					Version:        kv.Version,
				}
				if !arguments.KeyOnly {
					pair.Value = string(kv.Value)
				}
				response.Kvs = append(response.Kvs, pair)
			}

			return nil
		})
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Failed to get key: " + err.Error())), nil
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
