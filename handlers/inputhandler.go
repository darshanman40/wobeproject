package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/wobeproject/data"
	"github.com/wobeproject/util"
)

//InputHandler ...
func InputHandler(w http.ResponseWriter, r *http.Request) {
	var dataInput data.InputData

	r.ParseForm()
	if r.Method != "POST" {
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	_, readErr := buf.ReadFrom(r.Body)
	failOnError(readErr)
	body := buf.Bytes()
	err := json.Unmarshal(body, &dataInput)
	failOnError(err)

	newString := util.ReverseString(dataInput.Input)

	b, err := json.Marshal(data.InputData{
		Input: newString,
	})
	failOnError(err)
	w.Write(b)

}

func failOnError(err error) {
	if err != nil {
		log.Fatal("FAIL: ", err)
	}
}
