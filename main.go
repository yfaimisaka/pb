package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.POST("/", normalPaste)
	router.GET("/:pbid", contentByPbid)
	router.GET("/:pbid/:lan", contentByPbidHighLight)
	// router.ServeFiles("/home/aimi/workspace/go/pb/static/*filepath", http.Dir("/home/aimi/workspace/go/pb/static"))
	myLog.Fatal(http.ListenAndServe(":80", router))
}
