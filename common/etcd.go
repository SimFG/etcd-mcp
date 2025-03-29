package common

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func EtcdOp(connectionURL string, op func(context.Context, *clientv3.Client) error) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := clientv3.New(
		clientv3.Config{
			Endpoints: []string{connectionURL},
			Context:   ctx,
		},
	)
	if err != nil {
		return err
	}
	defer client.Close()

	return op(ctx, client)
}
