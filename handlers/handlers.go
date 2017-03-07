package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/wobeproject/logger"
)

var (
	l        logger.Logger
	httpErr  error
	postForm = []byte(`<html>
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
)

type page struct {
	title string
	body  []byte
}

//IndexHandler ...
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	l = logger.GetInstance()
	var inputHandleregex = regexp.MustCompile("^/$")
	httpErr = urlValidator(inputHandleregex, r.URL.Path)
	if httpErr != nil {
		return
	}
	contentType := r.Header.Get("Content-Type")
	l.Debug("Inside handler:", map[string]interface{}{
		"Handler id":   3,
		"Handler name": "inputHandler",
	})
	if r.Method != "POST" {
		p := &page{title: "Welcome", body: postForm}
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.title, p.body)
		return
	}

	rp := NewParser(w, contentType)
	if rp == nil {
		httpErr = errors.New("500: Internal Server error, Content-Type " + contentType + " not found")
		return
	}
	rp.RequestParse(r)
}

//ValidationHandler ...
func ValidationHandler(inner http.Handler) http.Handler {
	l = logger.GetInstance()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Debug("Inside handler:", map[string]interface{}{
			"Handler id":   2,
			"Handler name": "validationHandler",
		})

		defer invalidWriter(w)
		rMethod := r.Method

		l.Debug("Request", map[string]interface{}{
			"Method": rMethod,
		})
		if rMethod != "POST" {
			inner.ServeHTTP(w, r)
			return
		}
		contentType := r.Header.Get("Content-Type")
		l.Debug("Request Header", map[string]interface{}{
			"Content-Type": contentType,
		})

		if contentType == "" {
			httpErr = errors.New("400: Please provide Content-Type")
			return
		}

		if contentType != "application/x-www-form-urlencoded" {
			httpErr = errors.New("415: Invalid Content-Type, only accepts application/x-www-form-urlencoded")
			return
		}

		inner.ServeHTTP(w, r)
		l.Debug("exiting handler:", map[string]interface{}{
			"Handler id":   2,
			"Handler name": "validationHandler",
		})
	})
}

// RecoverHandler ....
func RecoverHandler(inner http.Handler) http.Handler {
	l = logger.GetInstance()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Debug("Inside handler:", map[string]interface{}{
			"Handler ID":   1,
			"Handler name": "recoverHandler",
		})
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
				l.Error("RECOVER HANDLER FAIL", map[string]interface{}{"ERR": err})
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		inner.ServeHTTP(w, r)
		l.Debug("exiting Handler", map[string]interface{}{
			"Handler id":   1,
			"Handler name": "recoverHandler",
		})
	})
}

//invalidWriter writes http error in response
func invalidWriter(w http.ResponseWriter) {
	l.Debug("Inside invalidWriter", map[string]interface{}{})

	if httpErr != nil {
		errMsg := httpErr.Error()
		status, _ := strconv.Atoi(errMsg[0:3])
		w.WriteHeader(status)
		w.Write([]byte(errMsg))
		l.Error("Inside validation defer", map[string]interface{}{
			"ERR": httpErr,
		})
	}
}

func urlValidator(handlerRegex *regexp.Regexp, path string) error {
	m := handlerRegex.FindStringSubmatch(path)
	if m == nil {
		return errors.New("404: No page found")
	}
	return nil
}
