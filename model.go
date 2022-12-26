package pb

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client


func init() {
    initRedisDB()
}

func initRedisDB() {
    rdb = redis.NewClient(&redis.Options{
        Addr: addr,
        Password: password,
        DB: 0, // use the default DB
    })

    if err := rdb.Ping(ctx).Err(); err != nil {
        panic(err)
    }
}






