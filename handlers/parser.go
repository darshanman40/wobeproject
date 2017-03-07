package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/wobeproject/logger"
	"github.com/wobeproject/util"
)

//Parser ...
type Parser interface {
	RequestParse(r *http.Request)
}

//revHandlerURLEncParser ...
type revHandlerURLEncParser struct {
	io.Writer
}

func (j *revHandlerURLEncParser) RequestParse(r *http.Request) {
	l.Debug("Content Type matched", map[string]interface{}{
		"Content-Type": r.Header.Get("Content-Type"),
	})

	r.ParseForm()
	input := r.Form.Get("input")
	if input == "" {
		httpErr = errors.New("422: No input found")
		return
	}

	newString := util.ReverseString(input)
	j.Write([]byte(newString))
}

//NewParser ...
func NewParser(w http.ResponseWriter, contentType string) Parser {
	l = logger.GetInstance()
	switch contentType {
	case "application/x-www-form-urlencoded":
		return &revHandlerURLEncParser{Writer: w}
	}
	return nil
}
