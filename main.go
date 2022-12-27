package main

import (
    "net/http"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.POST("/", normalPaste)
	router.GET("/:pbid", contentByPbid)
    myLog.Fatal(http.ListenAndServe(":8080", router))
}
