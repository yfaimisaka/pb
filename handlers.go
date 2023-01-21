package main

import (
	"errors"
	"fmt"
	"io"
    // use html/template instead of text/template to prevent xss..
	"html/template"

	// "io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type pbTmpl struct {
    Lan          string
	Theme        string
	PasteContent string
}

func normalPaste(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	var contents strings.Builder
	contents.WriteString(r.PostFormValue("c"))

	if contents.String() == "" {
		r.ParseMultipartForm(1<<63 - 2)
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
			"method":  "handlers.go: normalPaste",
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
	s = fmt.Sprintf(s, time.Now(), long, short, baseURL+short)
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
			"pbid":   pbid,
		},
	).Info("")

	// === serve for static resources like css, js ===
	// as the /static/xxx will conflict wiht /:pbid/:lan,
	// so use this way to solve static file (ugly though :<)
	// if slices.Contains[string](statics, pbid) {
	//     // LOG
	//     myLog.WithFields(
	//         logrus.Fields{
	//         "method": "handlers.go: contentByPbidHighLight",
	//            "resource": pbid,
	//     },).Info("try to get static resources")
	//     http.ServeFile(w, r, pbid)
	//     return
	// }

	if len(pbid) > 5 {
		var err1 error
		content, err1 = getV(pbid)
		err = errors.Join(err1)
	} else {
		var err2, err3 error
		long, err2 := getV(pbid)
		content, err3 = getV(long)
		err = errors.Join(err2, err3)
	}
	if err != nil {
		w.WriteHeader(http.StatusOK)
		myLog.Error(err)
	}

    // if find javascript
    if len(content) < 100 && strings.Contains(content, "script") {
        // set html to plain, prevent xss hack.
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    }

	rc := http.NewResponseController(w)
	rc.SetWriteDeadline(time.Time{})
	w.Write(stringToBytes(content))

}

func contentByPbidHighLight(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// pbid, may be short or long
	var content string
	var err error
	var theme = "github-dark"
	pbid := ps.ByName("pbid")
	lan := ps.ByName("lan")
	t := r.URL.Query().Get("t")
	if t != "" {
		theme = t
	}
	myLog.WithFields(
		logrus.Fields{
			"method": "handlers.go: contentByPbidHighLight",
			"pbid":   pbid,
			"lan":    lan,
			"theme":  theme,
		},
	).Info("")

	if len(pbid) > 5 {
		var err1 error
		content, err1 = getV(pbid)
		err = errors.Join(err1)
	} else {
		var err2, err3 error
		long, err2 := getV(pbid)
		content, err3 = getV(long)
		err = errors.Join(err2, err3)
	}

	pbtmp := pbTmpl{
		Theme:        theme,
		PasteContent: content,
	}
	// set default contentype
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// <Deprecated>
	// use template instead of pure text html.
	// rc := http.NewResponseController(w)
	// rc.SetWriteDeadline(time.Time{})
	// w.Write(stringToBytes(fmt.Sprintf(hlTemplate, theme, lan, content)))

	pbtmpl, err4 := template.New("pbtmpl.html").Delims("{[", "]}").ParseFiles("./pbtmpl.html")

    err = errors.Join(err4)

    err5 := pbtmpl.Execute(w, pbtmp)
	err = errors.Join(err5)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		myLog.Error(err)
	}
}

// === try to solve http force redirect ===
// func redirectToTls(w http.ResponseWriter, r *http.Request) {
//     http.Redirect(w, r, "https://aimisaka.site:443" + r.RequestURI, http.StatusMovedPermanently)
// }
