package common

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetEtcdClient(connectionURL string) (*clientv3.Client, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return clientv3.New(clientv3.Config{
		Endpoints: []string{connectionURL},
	})
}

func EtcdOp(connectionURL string, op func(*clientv3.Client) error) error {
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

	return op(client)
}
