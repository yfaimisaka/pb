package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func normalPaste(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var err error
    var contents strings.Builder
    r.ParseMultipartForm(1<<63-2)
    c := r.MultipartForm.Value["c"]
    for _, s := range c {
        contents.WriteString(s)
    }
    if contents.String() == "" {
        f := r.MultipartForm.File["c"]
        for _, fileContent := range f {
            fileContent, err0 := fileContent.Open()
            contentBytes, err1 := io.ReadAll(fileContent)
            err = errors.Join(err0, err1)
            contents.WriteString(bytesToString(contentBytes))
        }
    }
    content := contents.String()
    content = strings.TrimRight(content, "\n\r")
    // LOG
    myLog.WithField(
        "content", content,
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
    // pbid, may be short or long
    var content, short string
    var err1, err2, err3 error
	pbid := ps.ByName("pbid")
    myLog.WithField(
        "pbid", pbid,
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
	w.Write(stringToBytes(content))

}
