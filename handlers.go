package main

import (
	"errors"
	"fmt"
	"io"
	// "io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func normalPaste(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var err error
    var contents strings.Builder
    contents.WriteString(r.PostFormValue("c"))

    if contents.String() == "" {
        r.ParseMultipartForm(1<<63-2)
        // LOG
        myLog.WithFields(
            logrus.Fields{
                "method": "handlers.go: normalPaste",
            },
            ).Warn("cat not get 'c' in form-data value, try to get from form-data file")
        f := r.MultipartForm.File["c"]
        for _, fileContent := range f {
            fileContent, err0 := fileContent.Open()
            contentBytes, err1 := io.ReadAll(fileContent)
            err = errors.Join(err0, err1)
            contents.WriteString(bytesToString(contentBytes))
        }
        if contents.String() == "" {
            myLog.Warn("cat not get 'c' in form-data file, try to get from x-www-form-urlencoded")
            contents.WriteString(r.FormValue("c"))
        }
    }
    content := contents.String()
    content = strings.TrimRight(content, "\n\r")
    // LOG
    myLog.WithFields(
        logrus.Fields{
            "method": "handlers.go: normalPaste",
            "content": content,
        },
    ).Info("The content pasted")
    short, err2 := short(content)
    long, err3 := long(content)
    err4 := setV(long, short, content)
    err = errors.Join(err2, err3, err4)
    s := `
    time: %s
    long: %s
    short: %s
    url: %s
    `
    s = fmt.Sprintf(s, time.Now(), long, short, baseURL + short)
    if err != nil {
		w.WriteHeader(http.StatusOK)
		myLog.Error(err)
    } else {
        fmt.Fprintln(w, s)
    }
}

func contentByPbid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    var content string
    var err error
    // pbid, may be short or long
	pbid := ps.ByName("pbid")
    // LOG
    myLog.WithFields(
        logrus.Fields{
            "method": "handlers.go: contentByPbid",
            "pbid": pbid,
        },
    )
    if strings.Contains(pbid, "prism") {
        // LOG
        myLog.WithFields(
            logrus.Fields{
            "method": "handlers.go: contentByPbid",
        },).Info("try to get static resources")
        http.ServeFile(w, r, pbid)
        // w.Header().Write()
        // fbytes, err0 := ioutil.ReadFile(pbid)
        // w.WriteHeader(http.StatusOK)
        // w.Header().Set("Content-Type", "application/octet-stream")
        // w.Write(fbytes)
        // err = errors.Join(err0)
        return
    }

    if len(pbid) > 4 {
        var err1, err2 error
        short, err1 := getV(pbid)
        content, err2 = getV(short)
        err = errors.Join(err1, err2)
    } else {
        var err3 error
        content, err3 = getV(pbid)
        err = errors.Join(err3)
    }
	if err != nil {
		w.WriteHeader(http.StatusOK)
		myLog.Error(err)
	}

	rc := http.NewResponseController(w)
	rc.SetWriteDeadline(time.Time{})
	w.Write(stringToBytes(content))

}

func contentByPbidHighLight(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    // pbid, may be short or long
    var content, short string
    var err1, err2, err3 error
	pbid := ps.ByName("pbid")
    lan := ps.ByName("lan")
    myLog.WithFields(
        logrus.Fields{
            "pbid": pbid,
            "lan": lan,
        },
    ).Info("")
    if len(pbid) > 4 {
        short, err1 = getV(pbid)
        content, err2 = getV(short)
    } else {
        content, err3 = getV(pbid)
    }
    err := errors.Join(err1, err2, err3)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		myLog.Error(err)
	}

	rc := http.NewResponseController(w)
	rc.SetWriteDeadline(time.Time{})
	w.Write(stringToBytes(fmt.Sprintf(hlTemplate, lan, content)))
}

// func fileResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//     file := ps.ByName("file")  
//     contents, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		w.WriteHeader(http.StatusOK)
// 		myLog.Error(err)
// 	}
//     w.Write(contents)
// }
