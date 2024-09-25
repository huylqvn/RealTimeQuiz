package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedis(dsn string, isCluster bool) (client redis.UniversalClient, err error) {
	if isCluster {
		clusterOpts, err := redis.ParseClusterURL(dsn)
		if err != nil {
			return nil, err
		}
		client = redis.NewClusterClient(clusterOpts)
	} else {
		opts, err := redis.ParseURL(dsn)
		if err != nil {
			return nil, err
		}
		client = redis.NewClient(opts)
	}

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return client, err
}
