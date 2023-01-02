package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func main() {
	// === try to solve http force redirect ===
	// go func() {
	//     if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTls)); err != nil {
	//         myLog.WithFields(
	//             logrus.Fields{
	//                 "method": "main.go: main",
	//             },
	//         ).Error(err)
	//     }
	// }()
	router := httprouter.New()
	router.POST("/", normalPaste)
	router.GET("/:pbid", contentByPbid)
	router.GET("/:pbid/:lan", contentByPbidHighLight)
	// router.ServeFiles("/home/aimi/workspace/go/pb/static/*filepath", http.Dir("/home/aimi/workspace/go/pb/static"))

	// === try https === (failed)
	// go func() {
	//     router0 := httprouter.New()
	//     router0.POST("/", normalPaste)
	//     myLog.WithFields(
	//     logrus.Fields{
	//         "method": "main.go: main",
	//     },
	// ).Fatal(http.ListenAndServeTLS(":443", "aimisaka.site_bundle.crt", "aimisaka.site.key", router0))
	// }()

	myLog.WithFields(
		logrus.Fields{
			"method": "main.go: main",
		},
	).Fatal(http.ListenAndServe(":80", router))
}
