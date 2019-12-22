package dl

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type DelayQueue struct {
	*redis.Client
	alias string
}

func NewDelayQueue(redisClient *redis.Client, alias string) *DelayQueue {
	return &DelayQueue{
		redisClient, alias,
	}
}

func (q *DelayQueue) AddsDelay(values []interface{}, et time.Time) error {
	score := float64(et.Unix())
	delay := q.DelayName()

	pipe := q.Client.Pipeline()
	for i := range values {
		b, _ := json.Marshal(values[i])
		pipe.ZAdd(delay, redis.Z{
			Score:  score,
			Member: b,
		})
	}

	_, err := pipe.Exec()
	return err
}

func (q *DelayQueue) AddsQueue(values []interface{}) error {
	pipe := q.Client.Pipeline()
	queue := q.queueName()

	for i := range values {
		b, _ := json.Marshal(values[i])
		pipe.RPush(queue, b)
	}

	_, err := pipe.Exec()
	return err
}

func (q *DelayQueue) CheckAndSwap(n int64) (int, error) {
	count := 0
	for {
		rangeBy := redis.ZRangeBy{
			Min:    "-inf",
			Max:    strconv.Itoa(int(time.Now().Unix())),
			Offset: 0,
			Count:  n,
		}
		results, err := q.Client.ZRangeByScore(q.DelayName(), rangeBy).Result()
		if err != nil || len(results) == 0 {
			return count, err
		}

		err = q.swapQueue(results)
		if err != nil {
			return count, err
		}

		count += len(results)
	}

	return count, nil
}

func (q *DelayQueue) swapQueue(members []string) error {
	pipe := q.Client.Pipeline()
	pipe.ZRemRangeByRank(q.DelayName(), 0, int64(len(members)))
	queue := q.queueName()

	for i := range members {
		pipe.RPush(queue, members[i])
	}

	_, err := pipe.Exec()
	return err
}

func (q *DelayQueue) FetchQueue(n int64) ([]string, error) {
	results, err := q.Client.LRange(q.queueName(), 0, n-1).Result()
	if err != nil || len(results) == 0 {
		return nil, err
	}

	err = q.Client.LTrim(q.queueName(), int64(len(results)), -1).Err()
	if err != nil {
		pipe := q.Client.Pipeline()
		queue := q.queueName()

		for i := range results {
			pipe.RPush(queue, results[i])
		}

		pipe.Exec()
		return nil, err
	}

	return results, nil
}

func (q DelayQueue) Size() (int64, error) {
	return q.Client.ZCount(q.DelayName(), "-inf", "+inf").Result()
}

func (q DelayQueue) queueName() string {
	return "queue:" + q.alias
}

func (q DelayQueue) DelayName() string {
	return "delay:" + q.alias
}
