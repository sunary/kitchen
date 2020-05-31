package q

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// DelayQueue ...
type DelayQueue struct {
	*redis.Client
	QueueName string
	DelayName string
}

// NewDelayQueue ...
func NewDelayQueue(redisClient *redis.Client, alias string) *DelayQueue {
	return &DelayQueue{
		redisClient,
		"queue:" + alias,
		"delay:" + alias,
	}
}

// AddsDelay ...
func (q *DelayQueue) AddsDelay(values []interface{}, et time.Time) error {
	score := float64(et.Unix())
	members := make([]redis.Z, len(values))
	for i := range values {
		b, _ := json.Marshal(values[i])
		members[i] = redis.Z{
			Score:  score,
			Member: b,
		}
	}

	return q.Client.ZAdd(q.DelayName, members...).Err()
}

// AddsQueue ...
func (q *DelayQueue) AddsQueue(values []interface{}) error {
	members := make([]interface{}, len(values))
	for i := range values {
		members[i], _ = json.Marshal(values[i])
	}

	return q.Client.RPush(q.QueueName, members...).Err()
}

// CheckAndSwap ...
func (q *DelayQueue) CheckAndSwap(n int64) (int, error) {
	count := 0
	for {
		rangeBy := redis.ZRangeBy{
			Min:    "-inf",
			Max:    strconv.Itoa(int(time.Now().Unix())),
			Offset: 0,
			Count:  n,
		}
		results, err := q.Client.ZRangeByScore(q.DelayName, rangeBy).Result()
		if err != nil || len(results) == 0 {
			return count, err
		}

		err = q.swapQueue(results)
		if err != nil {
			return count, err
		}

		count += len(results)
	}
}

func (q *DelayQueue) swapQueue(members []string) error {
	pipe := q.Client.Pipeline()
	pipe.ZRemRangeByRank(q.DelayName, 0, int64(len(members)))

	mis := make([]interface{}, len(members))
	for i := range members {
		mis[i] = members[i]
	}

	pipe.RPush(q.QueueName, mis...)
	_, err := pipe.Exec()
	return err
}

// FetchQueue ...
func (q *DelayQueue) FetchQueue(n int64) ([]string, error) {
	results, err := q.Client.LRange(q.QueueName, 0, n-1).Result()
	if err != nil || len(results) == 0 {
		return nil, err
	}

	err = q.Client.LTrim(q.QueueName, int64(len(results)), -1).Err()
	if err != nil {
		members := make([]interface{}, len(results))
		for i := range results {
			members[i] = results[i]
		}

		q.Client.RPush(q.QueueName, members...)
		return nil, err
	}

	return results, nil
}

// Size ...
func (q DelayQueue) Size() (int64, error) {
	return q.Client.ZCount(q.DelayName, "-inf", "+inf").Result()
}
