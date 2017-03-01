package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wobeproject/data"
	"github.com/wobeproject/logger"
	"github.com/wobeproject/util"
)

var l = logger.GetInstance()

type page struct {
	title string
	body  []byte
}

//inputHandler ...
func inputHandler(w http.ResponseWriter, r *http.Request) {
	// var dataInput data.InputData

	if r.Method != "POST" {
		p := &page{title: "Welcome", body: []byte("Send POST request to reverse string")}
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.title, p.body)
		return
	}

	var b []byte
	r.ParseForm()

	switch r.Header.Get("Content-Type") {
	case "application/json":

		j := revHandlerJSONParser{Writer: w, d: data.InputData{}}
		j.requestParse(r)
		// l.Info("Content Type matched", map[string]interface{}{
		// 	"Content-Type": "application/json",
		// })
		// buf := bytes.NewBuffer(make([]byte, 0))
		// _, readErr := buf.ReadFrom(r.Body)
		// util.FailOnError("reading resp body failed ", readErr)
		// body := buf.Bytes()
		// err := json.Unmarshal(body, &dataInput)
		// util.FailOnError("json unmarshal failed ", err)
		// l.Info("Reversing string", map[string]interface{}{
		// 	"original string": dataInput.Input,
		// })
		// newString := util.ReverseString(dataInput.Input)
		//
		// b, err = json.Marshal(data.InputData{
		// 	Input: newString,
		// })
		// util.FailOnError("json marshal failed ", err)

	case "text/plain":

		l.Info("Content-Type matched: text/plain", map[string]interface{}{})

		v := r.FormValue("input")

		newString := util.ReverseString(v)

		b = []byte(newString)
		l.Info("Reversed string", map[string]interface{}{
			"newString": newString,
		})
	default:
		temp := r.FormValue("input")
		if temp == "" {
			l.Warning("No match for Content-Type Found in Request", map[string]interface{}{
				"Content-Type": r.Header.Get("Content-Type"),
			})
			b = []byte("ERROR No match for Content-Type Found in Request")
		} else {
			newString := util.ReverseString(temp)
			b = []byte(newString)
		}
	}
	w.Write(b)
}

// func validationHandler(inner http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Header.Get("Content-Type") {
// 		case "application/json":
//
// 		case "text/plain":
// 		default:
// 		}
//
// 		inner.ServeHTTP(w, r)
// 	})
// }

func recoverHandler(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		//	l := logger.GetInstance()
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				l.Panic("RECOVER HANDLER FAIL", map[string]interface{}{"ERR": err.Error()})
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		inner.ServeHTTP(w, r)
	})
}
