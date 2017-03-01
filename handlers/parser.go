package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/wobeproject/data"
	"github.com/wobeproject/util"
)

type revHandlerJSONParser struct {
	io.Writer
	d data.InputData
}

func (j *revHandlerJSONParser) requestParse(r *http.Request) {
	l.Info("Content Type matched", map[string]interface{}{
		"Content-Type": "application/json",
	})
	buf := bytes.NewBuffer(make([]byte, 0))
	_, readErr := buf.ReadFrom(r.Body)
	util.FailOnError("reading resp body failed ", readErr)
	body := buf.Bytes()
	err := json.Unmarshal(body, &j.d)
	util.FailOnError("json unmarshal failed ", err)
	l.Info("Reversing string", map[string]interface{}{
		"original string": j.d.Input,
	})
	newString := util.ReverseString(j.d.Input)

	b, err := json.Marshal(data.InputData{
		Input: newString,
	})
	//io.WriteString(j., s)
	// util.FailOnError("json marshal failed ", err)
	if b != nil {
		j.Write(b)
	} else {
		io.WriteString(j, err.Error())
	}
}
