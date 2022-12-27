package main

import (
	"errors"
	"testing"
)

func TestInitRedisDB(t *testing.T) {
	initRedisDB()

    var err error
    var shorts []string = make([]string, 0)

	contents := []string{
        "world",
         "world世界",
        `func main() { fmt.Println("hello world") }`,
	}

	for _, content := range contents {
        long, err0 := long(content)
        short, err1 := short(content)
        // append short to shorts
        shorts = append(shorts, short)
        err2 := setV(long, short, content)
        err = errors.Join(err0, err1, err2)
        if err != nil {
            t.Error(err)
        }
	}

    if len(shorts) != 3 {
        t.Errorf("got number of shorts: [%d], want: [3]", len(shorts))
    }

	for index, short := range shorts {
        content, err := getV(short);
		if err != nil {
			t.Error(err)
		}
        if content != contents[index] {
            t.Errorf("content of key: [%s], got: [%s], want: [%s]", short, content, contents[index])
        }

	}
}
