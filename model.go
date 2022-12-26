package pb

import (
	"context"
	"time"
	"unsafe"

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

func setV(key string, value []byte) (err error) {
    err = rdb.Set(ctx, key, unsafe.String(&value[0], len(value)), time.Hour * 24 * 5).Err()
    return
}

func getV(key string) (value string, err error) {
    strcmd := rdb.Get(ctx, key)
    value = strcmd.Val()
    err = strcmd.Err()
    return
}




