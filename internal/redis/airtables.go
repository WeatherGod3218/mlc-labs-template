package redis

import (
	"context"
	"math/rand/v2"

	"github.com/WeatherGod3218/mlc-project-template/internal/logging"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const QueueKey = "queued_tables"
const PointerKey = "queue_pointer"

func reshuffleQueue() error {
	ctx := context.Background()

	tables, err := client.LRange(ctx, QueueKey, 0, -1).Result()
	if err != nil {
		return err
	}

	rand.Shuffle(len(tables), func(i, j int) {
		tables[i], tables[j] = tables[j], tables[i]
	})

	_, err = client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Del(ctx, QueueKey)
		pipe.RPush(ctx, QueueKey, tables)
		pipe.Set(ctx, PointerKey, 0, 0)
		return nil
	})
	return err
}

func InitQueue(tables []string) error {
	ctx := context.Background()

	exists, err := client.Exists(ctx, QueueKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}

	rand.Shuffle(len(tables), func(i, j int) {
		tables[i], tables[j] = tables[j], tables[i]
	})

	_, err = client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.RPush(ctx, QueueKey, tables)
		pipe.Set(ctx, PointerKey, 0, 0)
		return nil
	})
	return err
}

func GetNextTable() (string, error) {
	ctx := context.Background()

	length, err := client.LLen(ctx, QueueKey).Result()
	if err != nil {
		return "", err
	}

	pointer, err := client.Get(ctx, PointerKey).Int64()
	if err != nil {
		return "", err
	}

	if pointer >= length {
		if err := reshuffleQueue(); err != nil {
			return "", err
		}
		pointer = 0
	}

	table, err := client.LIndex(ctx, QueueKey, pointer).Result()
	if err != nil {
		return "", err
	}

	err = client.Incr(ctx, PointerKey).Err()

	logging.Logger.WithFields(logrus.Fields{"table": table, "module": "api", "method": "GetData"}).Info("Selected New Table!")

	return table, err
}
