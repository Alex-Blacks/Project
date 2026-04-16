package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func GetBalance(id int) (int, error) {
	key := fmt.Sprintf("user:%d", id)

	if val, err := rdb.Get(ctx, key).Result(); err == nil {
		return strconv.Atoi(val)
	}

	balance := db[id]

	rdb.Set(ctx, key, balance, 10*time.Second)
	return balance, nil
}

func UpdateBalance(id int, amount int) error {
	db[id] = amount

	key := fmt.Sprintf("user:%d", id)
	rdb.Del(ctx, key)

	return nil
}
