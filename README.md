# etcd-mcp

etcd-mcp is a Management Control Panel server for etcd, providing a set of tools to interact with etcd clusters through a unified interface.

## Features

- Health Check: Monitor the health status of etcd clusters
- Get: Retrieve key-value pairs from etcd
- Put: Store key-value pairs in etcd
- Delete: Remove key-value pairs from etcd
- Status: Check etcd cluster status
- Save: Backup etcd data to files
- Restore: Restore etcd data from backup files

## Installation

```bash
go install github.com/SimFG/etcd-mcp
```

## Usage

The server provides the following tools:

- `health`: Check etcd cluster health
- `get`: Get key-value pairs from etcd
- `put`: Put key-value pairs into etcd
- `del`: Delete key-value pairs from etcd
- `status`: Get etcd cluster status
- `save`: Save etcd data to a backup file
- `restore`: Restore etcd data from a backup file

## Configuration

cursor config:

```
{
    "mcpServers": {
        "etcd": {
            "command": "[etcd-mcp bin file path]",
            "args": [],
            "env": {}
        }
    }
}
```

## Dependencies

- github.com/SimFG/etcd-mcp/tools
- github.com/metoro-io/mcp-golang
- github.com/metoro-io/mcp-golang/transport/stdio