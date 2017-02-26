package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/wobeproject/data"
	"github.com/wobeproject/logger"
	"github.com/wobeproject/util"
)

var l = logger.GetInstance()

//InputHandler ...
func InputHandler(w http.ResponseWriter, r *http.Request) {
	var dataInput data.InputData

	r.ParseForm()
	if r.Method != "POST" {
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	_, readErr := buf.ReadFrom(r.Body)
	failOnError("reading resp body failed ", readErr)
	body := buf.Bytes()
	err := json.Unmarshal(body, &dataInput)
	failOnError("json unmarshal failed ", err)

	newString := util.ReverseString(dataInput.Input)

	b, err := json.Marshal(data.InputData{
		Input: newString,
	})
	failOnError("json marshal failed ", err)
	w.Write(b)

}

func failOnError(msg string, err error) {
	if err != nil {
		l.Error(msg, map[string]interface{}{
			"ERR": err.Error(),
		})

	}
}
