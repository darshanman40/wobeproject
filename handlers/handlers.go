package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wobeproject/logger"
)

var l logger.Logger
var valid map[bool]string

type page struct {
	title string
	body  []byte
}

var postForm = []byte(`<html>
<head>
<title></title>
</head>
<body>
<form action="/" method="post">
    Text to reverse:  <input type="text" name="input"> <br > <br >
    <input type="submit" value="Input to Reverse">
</form>
</body>
</html>`)

//InputHandler ...
func InputHandler(w http.ResponseWriter, r *http.Request) {
	l = logger.GetInstance()
	l.Info("Inside handler:", map[string]interface{}{
		"Handler id":   3,
		"Handler name": "inputHandler",
	})
	if r.Method != "POST" {
		p := &page{title: "Welcome", body: postForm} //[]byte("Send POST request to reverse string")}
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.title, p.body)
		return
	}

	rp := NewParser(w, r.Header.Get("Content-Type"))
	if rp == nil {
		valid[false] = "500: Internal Server error"
		return
	}
	rp.RequestParse(r)

}

//ValidationHandler ...
func ValidationHandler(inner http.Handler) http.Handler {
	l = logger.GetInstance()
	l.Info("Inside handler:", map[string]interface{}{
		"Handler id":   2,
		"Handler name": "validationHandler",
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		valid = make(map[bool]string)
		defer invalidWriter(w, valid)
		if r.Method != "POST" {
			inner.ServeHTTP(w, r)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			valid[false] = "400: Please provide Content-Type"
			return
		}
		if contentType != "application/x-www-form-urlencoded" {
			valid[false] = "415: Invalid Content-Type, only accepts application/x-www-form-urlencoded"
			return
		}

		inner.ServeHTTP(w, r)
	})
}

func invalidWriter(w http.ResponseWriter, valid map[bool]string) {
	if msg := valid[false]; msg != "" {
		status, _ := strconv.Atoi(msg[0:3])
		w.WriteHeader(status)
		w.Write([]byte(msg))
		l.Error("Inside validation defer", map[string]interface{}{
			"valid[false]": valid[false],
			"valid[true]":  valid[true],
		})
	}
}

// RecoverHandler ....
func RecoverHandler(inner http.Handler) http.Handler {
	l = logger.GetInstance()
	l.Info("Inside handler:", map[string]interface{}{
		"Handler id":   1,
		"Handler name": "recoverHandler",
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
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
				l.Error("RECOVER HANDLER FAIL", map[string]interface{}{"ERR": err.Error()})
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		inner.ServeHTTP(w, r)
	})
}
