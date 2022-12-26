package pb

import (
	"testing"
)

func TestInitRedisDB(t *testing.T) {
    initRedisDB()

    kvs := map[string]string{
        "hello": "world", 
        "你好": "world世界",
        "AGX29mAYx7": `func main() { fmt.Println("hello world") }`,
    }
    for key, val := range kvs {
        if err := setV(key, []byte(val)); err != nil {
            t.Error(err)
        }
    }

    for key, _ := range kvs {
        if _, err := getV(key); err != nil {
            t.Error(err)
        } 
    }
}
